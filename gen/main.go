// Command gen emits thin Go accessor wrappers for the joltc C API functions
// that are not already hand-wrapped.
//
// It targets the "accessor" shape only:
//
//	void  JPH_T_GetX(const JPH_T*, Scalar* out)      -> func (t *T) X() Scalar
//	Scalar JPH_T_GetX(const JPH_T*)                  -> func (t *T) X() Scalar
//	void  JPH_T_SetX(JPH_T*, Scalar value)           -> func (t *T) SetX(value Scalar)
//	void  JPH_T_GetX(const JPH_T*, JPH_Vec3* out)    -> func (t *T) X() Vec3
//	...
//
// Vector/struct/pointer-returning and otherwise composite signatures are left
// to hand-written code. The generated file is committed.
//
// Usage: go run ./gen [-o internal/joltc/gen_accessors.go] [-check]
package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strings"
)

// goReceiver maps a C handle type to the Go type the generated method hangs off.
var goReceiver = map[string]string{
	"JPH_Body":                             "Body",
	"JPH_BodyCreationSettings":             "BodyCreationSettings",
	"JPH_BodyInterface":                    "BodyInterface",
	"JPH_PhysicsSystem":                    "PhysicsSystem",
	"JPH_Shape":                            "Shape",
	"JPH_JobSystem":                        "JobSystem",
	"JPH_BroadPhaseLayerInterface":         "BroadPhaseLayerInterface",
	"JPH_ObjectLayerPairFilter":            "ObjectLayerPairFilter",
	"JPH_ObjectVsBroadPhaseLayerFilter":    "ObjectVsBroadPhaseLayerFilter",
	"JPH_NarrowPhaseQuery":                 "NarrowPhaseQuery",
	"JPH_BroadPhaseQuery":                  "BroadPhaseQuery",
	"JPH_Skeleton":                         "Skeleton",
	"JPH_SkeletonPose":                     "SkeletonPose",
	"JPH_SkeletonMapper":                   "SkeletonMapper",
	"JPH_Ragdoll":                          "Ragdoll",
	"JPH_RagdollSettings":                  "RagdollSettings",
	"JPH_WheelSettings":                    "WheelSettings",
	"JPH_Wheel":                            "Wheel",
	"JPH_VehicleTransmissionSettings":      "VehicleTransmissionSettings",
	"JPH_WheeledVehicleControllerSettings": "WheeledVehicleControllerSettings",
	"JPH_Character":                        "Character",
	"JPH_CharacterBase":                    "Character",
	"JPH_Constraint":                       "Constraint",
}

// scalar maps a C scalar type to Go.
var scalar = map[string]struct{ goType, conv string }{
	"float":    {"float32", "float32"},
	"double":   {"float64", "float64"},
	"bool":     {"bool", "bool"},
	"int":      {"int", "int"},
	"int32_t":  {"int32", "int32"},
	"uint32_t": {"uint32", "uint32"},
	"int64_t":  {"int64", "int64"},
	"uint64_t": {"uint64", "uint64"},
	"uint8_t":  {"uint8", "uint8"},
	"size_t":   {"uint", "uint"},
}

// enums maps a C enum/typedef return type to a Go type defined in the package.
var enums = map[string]string{
	"JPH_PhysicsUpdateError":      "PhysicsUpdateError",
	"JPH_BodyType":                "BodyType",
	"JPH_MotionType":              "MotionType",
	"JPH_Activation":              "Activation",
	"JPH_ValidateResult":          "ValidateResult",
	"JPH_ShapeType":               "ShapeType",
	"JPH_ShapeSubType":            "ShapeSubType",
	"JPH_MotionQuality":           "MotionQuality",
	"JPH_OverrideMassProperties":  "OverrideMassProperties",
	"JPH_AllowedDOFs":             "AllowedDOFs",
	"JPH_GroundState":             "GroundState",
	"JPH_BackFaceMode":            "BackFaceMode",
	"JPH_ActiveEdgeMode":          "ActiveEdgeMode",
	"JPH_CollectFacesMode":        "CollectFacesMode",
	"JPH_CollisionCollectorType":  "CollisionCollectorType",
	"JPH_ConstraintType":          "ConstraintType",
	"JPH_ConstraintSubType":       "ConstraintSubType",
	"JPH_ConstraintSpace":         "ConstraintSpace",
	"JPH_MotorState":              "MotorState",
	"JPH_SpringMode":              "SpringMode",
	"JPH_SwingType":               "SwingType",
	"JPH_SixDOFConstraintAxis":    "SixDOFConstraintAxis",
	"JPH_Mesh_Shape_BuildQuality": "MeshShapeBuildQuality",
	"JPH_TransmissionMode":        "TransmissionMode",
	"JPH_BodyID":                  "BodyID",
	"JPH_SubShapeID":              "SubShapeID",
	"JPH_ObjectLayer":             "ObjectLayer",
	"JPH_BroadPhaseLayer":         "BroadPhaseLayer",
	"JPH_CollisionGroupID":        "CollisionGroupID",
	"JPH_CollisionSubGroupID":     "CollisionSubGroupID",
	"JPH_CharacterID":             "CharacterID",
}

// vectors maps a C vector/struct type to a Go type with converters.
var vectors = map[string]string{
	"JPH_Vec3":  "vecmath.Vec3",
	"JPH_Vec4":  "vecmath.Vec4",
	"JPH_Quat":  "vecmath.Quat",
	"JPH_Plane": "vecmath.Plane",
	"JPH_AABox": "vecmath.AABox",
	"JPH_RVec3": "vecmath.RVec3",
	"JPH_Mat4":  "vecmath.Mat4",
	"JPH_RMat4": "vecmath.RMat4",
}

var cArgs = map[string]string{
	"vecmath.Vec3":  "vec3ToC",
	"vecmath.Vec4":  "vec4ToC",
	"vecmath.Quat":  "quatToC",
	"vecmath.Plane": "planeToC",
	"vecmath.AABox": "aaboxToC",
	"vecmath.RVec3": "rvec3ToC",
	"vecmath.Mat4":  "mat4ToC",
	"vecmath.RMat4": "rmat4ToC",
}

var cFrom = map[string]string{
	"vecmath.Vec3":  "vec3FromC",
	"vecmath.Vec4":  "vec4FromC",
	"vecmath.Quat":  "quatFromC",
	"vecmath.Plane": "planeFromC",
	"vecmath.AABox": "aaboxFromC",
	"vecmath.RVec3": "rvec3FromC",
	"vecmath.Mat4":  "mat4FromC",
	"vecmath.RMat4": "rmat4FromC",
}

type param struct {
	cType string // normalized, e.g. "JPH_Vec3", "float", "JPH_Body*"
	ptr   bool
	cst   bool
}

func (p param) base() string { return p.cType }

func normalize(t string) param {
	t = strings.TrimSpace(t)
	t = strings.ReplaceAll(t, "JPH_API_CALL", "")
	t = strings.Join(strings.Fields(t), " ")
	// Drop the trailing parameter name: keep everything up to the last '*', or
	// drop the final token when there is no pointer.
	if i := strings.LastIndex(t, "*"); i >= 0 {
		t = t[:i+1]
	} else if i := strings.LastIndex(t, " "); i >= 0 {
		t = t[:i]
	}
	t = strings.TrimSpace(t)
	p := param{}
	if strings.HasPrefix(t, "const ") {
		p.cst = true
		t = strings.TrimPrefix(t, "const ")
	}
	if strings.HasSuffix(t, "*") {
		p.ptr = true
		t = strings.TrimSpace(strings.TrimSuffix(t, "*"))
	}
	p.cType = t
	return p
}

func main() {
	out := flag.String("o", "internal/joltc/gen_accessors.go", "output file")
	header := flag.String("header", "third_party/joltc/include/joltc.h", "joltc header")
	check := flag.Bool("check", false, "verify the output file is up to date instead of writing it")
	flag.Parse()

	decls, err := parseHeader(*header)
	if err != nil {
		fatal(err)
	}
	wrapped := wrappedFunctions("internal/joltc")

	var b strings.Builder
	b.WriteString("//go:build cgo\n\n")
	b.WriteString("// Code generated by go run ./gen; DO NOT EDIT.\n\n")
	b.WriteString("package joltc\n\n")
	b.WriteString("/*\n#include <joltc.h>\n*/\nimport \"C\"\n\n")
	b.WriteString("import (\n\t\"unsafe\"\n\n\t\"github.com/akhenakh/gojoltc/internal/vecmath\"\n)\n\n")
	b.WriteString("// Generated accessors for joltc functions not covered by the hand-written API.\n\n")
	b.WriteString("var _ = vecmath.Vec3{}\n\n")

	generated, skipped := 0, 0
	generatedNames := map[string]bool{}
	var skippedReasons []string
	for _, d := range decls {
		if wrapped[d.name] {
			continue
		}
		code, reason := emit(d)
		if code == "" {
			skipped++
			if reason != "" && len(skippedReasons) < 40 {
				skippedReasons = append(skippedReasons, fmt.Sprintf("%s: %s", d.name, reason))
			}
			continue
		}
		b.WriteString(code)
		generated++
		generatedNames[d.name] = true
	}

	src := b.String()
	if formatted, err := format.Source([]byte(src)); err == nil {
		src = string(formatted)
	} else {
		fmt.Fprintln(os.Stderr, "gen: warning: generated source is not gofmt-clean:", err)
	}

	if *check {
		existing, err := os.ReadFile(*out)
		if err != nil || string(existing) != src {
			fmt.Fprintf(os.Stderr, "gen: %s is out of date; run: go run ./gen\n", *out)
			os.Exit(1)
		}
		fmt.Printf("generated %d accessors, skipped %d; %s is up to date\n", generated, skipped, *out)
		return
	}

	if err := os.WriteFile(*out, []byte(src), 0o644); err != nil {
		fatal(err)
	}
	if err := writeCoverage(decls, wrapped, generatedNames); err != nil {
		fatal(err)
	}
	fmt.Printf("generated %d accessors, skipped %d\n", generated, skipped)
	if len(skippedReasons) > 0 {
		fmt.Println("sample skips:")
		for _, r := range skippedReasons {
			fmt.Println("  -", r)
		}
	}
}

type decl struct {
	ret    param
	name   string
	params []param
}

var declRe = regexp.MustCompile(`(?s)^JPH_CAPI\s+(.*?)\s+(JPH_[A-Za-z0-9_]+)\s*\((.*)\)\s*$`)

func parseHeader(path string) ([]decl, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Strip block comments, then split declarations on ';'.
	src := regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(string(raw), "")
	var out []decl
	for part := range strings.SplitSeq(src, ";") {
		part = strings.TrimSpace(part)
		if !strings.HasPrefix(part, "JPH_CAPI") {
			continue
		}
		part = strings.Join(strings.Fields(part), " ")
		m := declRe.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		d := decl{ret: normalize(m[1]), name: m[2]}
		if params := strings.TrimSpace(m[3]); params != "" && params != "void" {
			for p := range strings.SplitSeq(params, ",") {
				d.params = append(d.params, normalize(p))
			}
		}
		out = append(out, d)
	}
	return out, nil
}

var refRe = regexp.MustCompile(`C\.(JPH_[A-Za-z0-9_]+)`)

// wrappedFunctions returns the set of C function names already referenced from
// the hand-written cgo sources.
func wrappedFunctions(dir string) map[string]bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fatal(err)
	}
	set := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasPrefix(e.Name(), "gen_") {
			continue
		}
		data, err := os.ReadFile(dir + "/" + e.Name())
		if err != nil {
			fatal(err)
		}
		for _, m := range refRe.FindAllStringSubmatch(string(data), -1) {
			set[m[1]] = true
		}
	}
	return set
}

// emit returns generated Go code for an accessor-shaped declaration, or "" with
// a reason when the shape is unsupported.
func emit(d decl) (code, reason string) {
	parts := strings.SplitN(strings.TrimPrefix(d.name, "JPH_"), "_", 2)
	if len(parts) != 2 {
		return "", "no type prefix"
	}
	cType := "JPH_" + parts[0]
	recv, ok := goReceiver[cType]
	if !ok {
		return "", "no Go receiver for " + cType
	}
	if len(d.params) == 0 || !d.params[0].ptr || d.params[0].base() != cType {
		return "", "unexpected receiver"
	}
	goName := methodName(parts[1])
	if goName == "" {
		return "", "unnamed method"
	}
	// joltc uses different C struct types for base/derived handles (e.g.
	// JPH_Character vs JPH_CharacterBase); reinterpret the receiver pointer.
	rExpr := fmt.Sprintf("(*C.%s)(unsafe.Pointer(r.ptr))", cType)

	// Optional leading handle id (BodyInterface and friends take one).
	rem := d.params[1:]
	idDecl, idArgs := "", ""
	if len(rem) > 0 && !rem[0].ptr {
		if gt, ok := enums[rem[0].base()]; ok && (rem[0].base() == "JPH_BodyID" || rem[0].base() == "JPH_CharacterID") {
			idDecl = "id " + gt + ", "
			idArgs = fmt.Sprintf("C.%s(id)", rem[0].base())
			rem = rem[1:]
		}
	}

	isVoid := d.ret.base() == "void" && !d.ret.cst && !d.ret.ptr

	// Setter: void NAME(recv, [id,] T value) or const T* value.
	if isVoid && len(rem) == 1 {
		arg := rem[0]
		if gt, ok := scalar[arg.base()]; ok && !arg.ptr {
			return setter(recv, goName, d.name, rExpr, idDecl, idArgs, arg.base(), gt.goType), ""
		}
		if arg.ptr && arg.cst {
			if gt, ok := vectors[arg.base()]; ok {
				return vecSetter(recv, goName, d.name, rExpr, idDecl, idArgs, gt), ""
			}
		}
		if gt, ok := enums[arg.base()]; ok && !arg.ptr {
			return setter(recv, goName, d.name, rExpr, idDecl, idArgs, arg.base(), gt), ""
		}
	}

	// Value return: T NAME(recv, [id]).
	if !d.ret.ptr {
		if len(rem) == 0 {
			if gt, ok := scalar[d.ret.base()]; ok {
				return retScalar(recv, goName, d.name, rExpr, idDecl, idArgs, gt.goType, gt.conv), ""
			}
			if gt, ok := enums[d.ret.base()]; ok {
				return retEnum(recv, goName, d.name, rExpr, idDecl, idArgs, gt), ""
			}
		}
		if len(rem) == 0 && d.ret.cst {
			if gt, ok := vectors[d.ret.base()]; ok {
				return retVec(recv, goName, d.name, rExpr, idDecl, idArgs, gt), ""
			}
		}
	}

	// Out-parameter: void NAME(recv, [id,] T* out).
	if isVoid && len(rem) == 1 && rem[0].ptr && !rem[0].cst {
		out := rem[0].base()
		if gt, ok := scalar[out]; ok {
			return outScalar(recv, goName, d.name, rExpr, idDecl, idArgs, gt.goType, gt.conv), ""
		}
		if gt, ok := enums[out]; ok {
			return outEnum(recv, goName, d.name, rExpr, idDecl, idArgs, gt), ""
		}
		if gt, ok := vectors[out]; ok {
			return outVec(recv, goName, d.name, rExpr, idDecl, idArgs, gt), ""
		}
	}
	return "", "unsupported shape"
}

func methodName(m string) string {
	m = strings.TrimPrefix(m, "Get")
	if m == "" {
		return ""
	}
	return m
}

// joinArgs joins the non-empty parts of a C argument list.
func joinArgs(parts ...string) string {
	var out []string
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ", ")
}

func setter(recv, goName, cName, rExpr, idDecl, idArgs, cArg, goType string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%svalue %s) {\n\tC.%s(%s)\n}\n\n",
		goName, cName, recv, goName, idDecl, goType, cName, joinArgs(rExpr, idArgs, fmt.Sprintf("C.%s(value)", cArg)))
}

func vecSetter(recv, goName, cName, rExpr, idDecl, idArgs, goType string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%svalue %s) {\n\tc := %s(value)\n\tC.%s(%s)\n}\n\n",
		goName, cName, recv, goName, idDecl, goType, cArgs[goType], cName, joinArgs(rExpr, idArgs, "&c"))
}

func retScalar(recv, goName, cName, rExpr, idDecl, idArgs, goType, conv string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%s) %s {\n\treturn %s(C.%s(%s))\n}\n\n",
		goName, cName, recv, goName, strings.TrimSuffix(idDecl, ", "), goType, conv, cName, joinArgs(rExpr, idArgs))
}

func retEnum(recv, goName, cName, rExpr, idDecl, idArgs, goType string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%s) %s {\n\treturn %s(C.%s(%s))\n}\n\n",
		goName, cName, recv, goName, strings.TrimSuffix(idDecl, ", "), goType, goType, cName, joinArgs(rExpr, idArgs))
}

func retVec(recv, goName, cName, rExpr, idDecl, idArgs, goType string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%s) %s {\n\tc := C.%s(%s)\n\treturn %s(*c)\n}\n\n",
		goName, cName, recv, goName, strings.TrimSuffix(idDecl, ", "), goType, cName, joinArgs(rExpr, idArgs), cFrom[goType])
}

func outScalar(recv, goName, cName, rExpr, idDecl, idArgs, goType, conv string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%s) %s {\n\tvar out C.%s\n\tC.%s(%s)\n\treturn %s(out)\n}\n\n",
		goName, cName, recv, goName, strings.TrimSuffix(idDecl, ", "), goType, cScalar(goType), cName, joinArgs(rExpr, idArgs, "&out"), conv)
}

func outEnum(recv, goName, cName, rExpr, idDecl, idArgs, goType string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%s) %s {\n\tvar out C.%s\n\tC.%s(%s)\n\treturn %s(out)\n}\n\n",
		goName, cName, recv, goName, strings.TrimSuffix(idDecl, ", "), goType, cType(goType), cName, joinArgs(rExpr, idArgs, "&out"), goType)
}

func outVec(recv, goName, cName, rExpr, idDecl, idArgs, goType string) string {
	return fmt.Sprintf("// %s wraps %s.\nfunc (r *%s) %s(%s) %s {\n\tvar out C.%s\n\tC.%s(%s)\n\treturn %s(out)\n}\n\n",
		goName, cName, recv, goName, strings.TrimSuffix(idDecl, ", "), goType, cType(goType), cName, joinArgs(rExpr, idArgs, "&out"), cFrom[goType])
}

func cType(goType string) string {
	for c, g := range enums {
		if g == goType {
			return c
		}
	}
	for c, g := range vectors {
		if g == goType {
			return c
		}
	}
	return goType
}

func cScalar(goType string) string {
	for c, s := range scalar {
		if s.goType == goType {
			return c
		}
	}
	return goType
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "gen:", err)
	os.Exit(1)
}

// writeCoverage writes docs/coverage.md summarising how much of the joltc C API
// is wrapped, and which groups remain.
func writeCoverage(decls []decl, wrapped map[string]bool, generated map[string]bool) error {
	handWrapped := 0
	remaining := map[string]int{}
	for _, d := range decls {
		switch {
		case wrapped[d.name]:
			handWrapped++
		case generated[d.name]:
		default:
			g := d.name
			parts := strings.SplitN(strings.TrimPrefix(d.name, "JPH_"), "_", 2)
			if len(parts) == 2 {
				g = "JPH_" + parts[0]
			}
			remaining[g]++
		}
	}

	type kv struct {
		k string
		v int
	}
	var groups []kv
	for k, v := range remaining {
		groups = append(groups, kv{k, v})
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].v != groups[j].v {
			return groups[i].v > groups[j].v
		}
		return groups[i].k < groups[j].k
	})

	var b strings.Builder
	b.WriteString("# joltc API coverage\n\n")
	b.WriteString("Generated by `go run ./gen`; do not edit by hand.\n\n")
	fmt.Fprintf(&b, "- header functions parsed: **%d**\n", len(decls))
	fmt.Fprintf(&b, "- hand-wrapped: **%d**\n", handWrapped)
	fmt.Fprintf(&b, "- generated accessors: **%d**\n", len(generated))
	total := handWrapped + len(generated)
	fmt.Fprintf(&b, "- wrapped total: **%d / %d**\n\n", total, len(decls))
	b.WriteString("The remaining functions are composite (create/destroy taking or returning pointers to\n")
	b.WriteString("POD structs), callbacks, or the intentionally skipped C math helpers (`JPH_Vec3_*`,\n")
	b.WriteString("`JPH_Quat_*`, `JPH_Mat4_*` — math is implemented in pure Go). They are wrapped by hand\n")
	b.WriteString("as needed rather than generated.\n\n")
	b.WriteString("## Remaining by group\n\n| group | functions |\n| --- | ---: |\n")
	for _, g := range groups {
		fmt.Fprintf(&b, "| `%s` | %d |\n", g.k, g.v)
	}
	return os.WriteFile("docs/coverage.md", []byte(b.String()), 0o644)
}
