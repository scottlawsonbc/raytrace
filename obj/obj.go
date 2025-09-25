// Package obj parses Wavefront OBJ files, a popular plain text 3D file format.
//
// It supports parsing vertices, texture coordinates, normals, faces, and
// materials defined in accompanying MTL files.
// It returns parsed data in a structured format representing the input OBJ file.
// Other packages can reference the parsed data to render, transform, or convert
// to a different format, such as a node for a particular rendering engine.
//
// The problem of loading OBJ files is somewhat analogous to the problem of
// template parsing, such as in `html/template` or `text/template`.
// These packages implement `html/template/parse` and `text/template/parse`
// which has loosely inspired the design of this package.
// Both template parsing and OBJ parsing involve reading text files, parsing
// them and constructing a structured representation of the data, combining
// data from multiple files (e.g., OBJ and MTL files) or sources (e.g., template
// text and data).
//
// # OBJ Files
//
// OBJ files are ASCII files that define 3D geometry, including vertices,
// faces, normals, and texture coordinates. There is no compression, and the
// file contains 1 or more named 3D objects.
//
// # MTL Files
//
// An MTL file is an auxiliary file that contains definitions of materials referenced by an OBJ file.
// The OBJ file specifies the MTL file using a directive such as:
//
//	mtllib file_name.mtl
//
// The MTL file defines various materials, for example, "shinyred" or "iron". Within the OBJ file,
// the directive
//
//	usemtl shinyred
//
// indicates that all subsequent faces should be rendered with this material until a new material is invoked.
//
// An MTL file consists of a sequence of material definitions. Each definition starts with a `newmtl` statement,
// followed by lines specifying the material's properties.
//
// Example MTL File:
//
//	newmtl shinyred
//	Ka  0.1986  0.0000  0.0000
//	Kd  0.5922  0.0166  0.0000
//	Ks  0.5974  0.2084  0.2084
//	illum 2
//	Ns 100.2237
//
// Characteristics of MTL Files:
//
// - **Format**: ASCII
// - **Comments**: Begin with a `#` character in column 1.
// - **Structure**: Consists of a sequence of `newmtl` statements, each defining a new material.
// - **Flexibility**: Blank lines may be inserted for clarity.
//
// Material Properties:
//
// Each material definition can include the following properties:
//
// - **Ka r g b**
//   - *Description*: Defines the ambient color of the material as (r, g, b).
//   - *Default*: (0.2, 0.2, 0.2)
//
// - **Kd r g b**
//   - *Description*: Defines the diffuse color of the material as (r, g, b).
//   - *Default*: (0.8, 0.8, 0.8)
//
// - **Ks r g b**
//   - *Description*: Defines the specular color of the material as (r, g, b), which appears in highlights.
//   - *Default*: (1.0, 1.0, 1.0)
//
// - **d alpha**
//   - *Description*: Defines the non-transparency of the material. `alpha` ranges from 0.0 (fully transparent) to 1.0 (fully opaque).
//   - *Default*: 1.0 (not transparent)
//
// - **Tr alpha**
//   - *Description*: Defines the transparency of the material. `alpha` ranges from 0.0 (not transparent) to 1.0 (fully transparent).
//   - *Default*: 0.0 (not transparent)
//   - *Note*: `d` and `Tr` are inversely related; specifying one affects the other.
//
// - **Ns s**
//   - *Description*: Defines the shininess of the material, where higher values result in smaller, sharper highlights.
//   - *Default*: 0.0
//
// - **illum n**
//   - *Description*: Specifies the illumination model used by the material.
//   - `illum = 1`: Flat material with no specular highlights; `Ks` is not used.
//   - `illum = 2`: Material has specular highlights; `Ks` must be specified.
//   - *Default*: Varies based on implementation.
//
// - **map_Ka filename**
//   - *Description*: Specifies a texture map file for the ambient color. The file should contain an ASCII dump of RGB values.
//   - *Default*: None

package obj

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strconv"
	"strings"
)

// Vertex represents a 3D point in space.
type Vertex struct {
	X, Y, Z float64
}

// TexCoord represents a 2D texture coordinate.
type TexCoord struct {
	U, V float64
}

// Normal represents a 3D normal vector.
type Normal struct {
	X, Y, Z float64
}

// Index holds 1-based indices to vertex data, referencing positions, texture coordinates, and normals.
type Index struct {
	Vertex   int // Index into the Vertices slice (1-based)
	TexCoord int // Index into the TexCoords slice (optional, 0 if not specified)
	Normal   int // Index into the Normals slice (optional, 0 if not specified)
}

// Face represents a polygonal face, defined by a list of indices to vertex data.
type Face struct {
	Indices  []Index // Indices defining the vertices of the face
	Material string  // Name of the material applied to the face (optional)
}

// Material defines the properties of a material, as specified in an MTL file.
type Material struct {
	Name      string     // Material name
	Diffuse   [3]float64 // Diffuse color (Kd)
	Ambient   [3]float64 // Ambient color (Ka)
	Specular  [3]float64 // Specular color (Ks)
	Shininess float64    // Specular exponent (Ns)
	Texture   string     // Texture filename (map_Kd)
}

// Object represents the contents of an OBJ file, including geometry and materials.
type Object struct {
	Vertices  []Vertex             // List of vertices (positions)
	TexCoords []TexCoord           // List of texture coordinates
	Normals   []Normal             // List of normals
	Faces     []Face               // List of faces defining the geometry
	Materials map[string]*Material // Map of material names to their definitions
}

// ParseError represents a parsing error with contextual information.
type ParseError struct {
	Filename string // Name of the file being parsed
	Line     int    // Line number where the error occurred
	LineText string // Content of the line where the error occurred
	Msg      string // Description of the error
}

// Error implements the error interface for ParseError.
func (e *ParseError) Error() string {
	if e.Filename != "" {
		return fmt.Sprintf("%s:%d: %s\n    %s", e.Filename, e.Line, e.Msg, e.LineText)
	}
	return fmt.Sprintf("line %d: %s\n    %s", e.Line, e.Msg, e.LineText)
}

// ParseFS reads and parses an OBJ file from the provided filesystem using the given pattern.
// It returns an Object containing the parsed geometry and material information.
func ParseFS(fsys fs.FS, pattern string) (*Object, error) {
	data, err := fs.ReadFile(fsys, path.Base(pattern))
	if err != nil {
		return nil, &ParseError{
			Filename: path.Base(pattern),
			Line:     0,
			LineText: "",
			Msg:      fmt.Sprintf("failed to read file '%s': %v", pattern, err),
		}
	}
	p := &parser{
		reader:   bufio.NewReader(bytes.NewReader(data)),
		obj:      &Object{Materials: make(map[string]*Material)},
		fsys:     fsys,
		filename: path.Base(pattern),
	}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p.obj, nil
}

// parser encapsulates the parsing state and logic.
type parser struct {
	reader          *bufio.Reader // Reader to read the OBJ file line by line
	obj             *Object       // Object being constructed
	currentMaterial string        // Current material name in use
	lineNumber      int           // Current line number in the OBJ file
	lineText        string        // Content of the current line
	fsys            fs.FS         // Filesystem to load external resources (e.g., MTL files)
	filename        string        // Name of the OBJ file being parsed
}

// parse initiates the parsing process of the OBJ file.
func (p *parser) parse() error {
	for {
		line, err := p.reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return &ParseError{
				Filename: p.filename,
				Line:     p.lineNumber,
				LineText: "",
				Msg:      fmt.Sprintf("error reading OBJ data: %v", err),
			}
		}
		// Handle the last line if it doesn't end with '\n'
		if err == io.EOF && len(line) == 0 {
			break
		}
		p.lineNumber++
		p.lineText = strings.TrimSpace(line)
		if err := p.parseLine(p.lineText); err != nil {
			return err
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

// parseLine processes a single line of the OBJ file.
func (p *parser) parseLine(line string) error {
	if line == "" || strings.HasPrefix(line, "#") {
		// Skip empty lines and comments.
		return nil
	}
	// Find the first space to determine the directive
	firstSpace := strings.IndexByte(line, ' ')
	if firstSpace == -1 {
		// Line has only one token, possibly invalid
		return nil // Or handle single-token directives if any
	}
	directive := line[:firstSpace]
	rest := line[firstSpace+1:]

	switch directive {
	case "v":
		return p.parseVertex(rest)
	case "vt":
		return p.parseTexCoord(rest)
	case "vn":
		return p.parseNormal(rest)
	case "f":
		return p.parseFace(rest)
	case "mtllib":
		return p.parseMTLLib(rest)
	case "usemtl":
		return p.parseUseMTL(rest)
	default:
		// Ignore unrecognized or unsupported directives
		return nil
	}
}

// parseVertex parses a vertex (position) definition.
func (p *parser) parseVertex(rest string) error {
	// Expecting three float values separated by spaces
	fields := splitFields(rest, 3)
	if len(fields) < 3 {
		return p.newError("invalid vertex data: expected at least 3 components, got %d", len(fields))
	}
	x, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return p.newError("invalid vertex X coordinate: %v", err)
	}
	y, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return p.newError("invalid vertex Y coordinate: %v", err)
	}
	z, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return p.newError("invalid vertex Z coordinate: %v", err)
	}
	p.obj.Vertices = append(p.obj.Vertices, Vertex{X: x, Y: y, Z: z})
	return nil
}

// parseTexCoord parses a texture coordinate definition.
func (p *parser) parseTexCoord(rest string) error {
	// Expecting two float values separated by spaces
	fields := splitFields(rest, 2)
	if len(fields) < 2 {
		return p.newError("invalid texture coordinate data: expected at least 2 components, got %d", len(fields))
	}
	u, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return p.newError("invalid texture U coordinate: %v", err)
	}
	v, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return p.newError("invalid texture V coordinate: %v", err)
	}
	p.obj.TexCoords = append(p.obj.TexCoords, TexCoord{U: u, V: v})
	return nil
}

// parseNormal parses a normal vector definition.
func (p *parser) parseNormal(rest string) error {
	// Expecting three float values separated by spaces
	fields := splitFields(rest, 3)
	if len(fields) < 3 {
		return p.newError("invalid normal data: expected at least 3 components, got %d", len(fields))
	}
	x, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return p.newError("invalid normal X component: %v", err)
	}
	y, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return p.newError("invalid normal Y component: %v", err)
	}
	z, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return p.newError("invalid normal Z component: %v", err)
	}
	p.obj.Normals = append(p.obj.Normals, Normal{X: x, Y: y, Z: z})
	return nil
}

// parseFace parses a face definition, which can reference vertices, texture coordinates, and normals.
func (p *parser) parseFace(rest string) error {
	// Faces can have varying number of vertices (usually 3 or 4)
	// Each vertex can have the format v, v/vt, v//vn, or v/vt/vn
	parts := splitFields(rest, -1) // Get all parts
	if len(parts) < 3 {
		return p.newError("face definition error: a face must have at least 3 vertices, got %d", len(parts))
	}
	var indices []Index
	indices = make([]Index, 0, len(parts)) // Preallocate with the number of vertices
	for _, part := range parts {
		index, err := p.parseIndex(part)
		if err != nil {
			return p.newError("invalid face index '%s': %v", part, err)
		}
		indices = append(indices, index)
	}
	p.obj.Faces = append(p.obj.Faces, Face{
		Indices:  indices,
		Material: p.currentMaterial,
	})
	return nil
}

// parseIndex parses a vertex reference in a face, which may include vertex, texture coordinate, and normal indices.
func (p *parser) parseIndex(s string) (Index, error) {
	var idx Index
	if s == "" {
		return idx, fmt.Errorf("empty face index")
	}
	parts := strings.Split(s, "/")
	switch len(parts) {
	case 1:
		// v
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return idx, fmt.Errorf("invalid vertex index: %v", err)
		}
		idx.Vertex, err = resolveIndex(v, len(p.obj.Vertices))
		if err != nil {
			return idx, err
		}
	case 2:
		// v/vt
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return idx, fmt.Errorf("invalid vertex index: %v", err)
		}
		vt, err := strconv.Atoi(parts[1])
		if err != nil {
			return idx, fmt.Errorf("invalid texture coordinate index: %v", err)
		}
		idx.Vertex, err = resolveIndex(v, len(p.obj.Vertices))
		if err != nil {
			return idx, err
		}
		idx.TexCoord, err = resolveIndex(vt, len(p.obj.TexCoords))
		if err != nil {
			return idx, err
		}
	case 3:
		// v//vn or v/vt/vn
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return idx, fmt.Errorf("invalid vertex index: %v", err)
		}
		idx.Vertex, err = resolveIndex(v, len(p.obj.Vertices))
		if err != nil {
			return idx, err
		}
		if parts[1] != "" {
			vt, err := strconv.Atoi(parts[1])
			if err != nil {
				return idx, fmt.Errorf("invalid texture coordinate index: %v", err)
			}
			idx.TexCoord, err = resolveIndex(vt, len(p.obj.TexCoords))
			if err != nil {
				return idx, err
			}
		}
		if parts[2] != "" {
			vn, err := strconv.Atoi(parts[2])
			if err != nil {
				return idx, fmt.Errorf("invalid normal index: %v", err)
			}
			idx.Normal, err = resolveIndex(vn, len(p.obj.Normals))
			if err != nil {
				return idx, err
			}
		}
	default:
		return idx, fmt.Errorf("invalid face index format")
	}
	return idx, nil
}

// resolveIndex resolves negative indices and checks range.
func resolveIndex(val, size int) (int, error) {
	if val < 0 {
		val = size + val + 1
	}
	if val < 1 || val > size {
		return 0, fmt.Errorf("index %d out of range (1 to %d)", val, size)
	}
	return val, nil
}

// parseUseMTL handles the usemtl directive, setting the current material for subsequent faces.
func (p *parser) parseUseMTL(rest string) error {
	if rest == "" {
		return p.newError("usemtl directive error: material name is missing")
	}
	p.currentMaterial = rest
	return nil
}

// parseMTLLib handles the mtllib directive, loading material definitions from an external MTL file.
func (p *parser) parseMTLLib(rest string) error {
	if rest == "" {
		return p.newError("mtllib directive error: filename is missing")
	}
	// Handle multiple mtllib directives by splitting filenames
	filenames := strings.Fields(rest)
	for _, filename := range filenames {
		baseFilename := path.Base(filename)
		data, err := fs.ReadFile(p.fsys, baseFilename)
		if err != nil {
			return p.newError("failed to read material library '%s': %v", baseFilename, err)
		}
		if err := p.parseMTL(bytes.NewReader(data), baseFilename); err != nil {
			return err
		}
	}
	return nil
}

// parseMTL parses an MTL file, populating the Materials map with material definitions.
func (p *parser) parseMTL(r io.Reader, mtlFilename string) error {
	scanner := bufio.NewScanner(r)
	// Increase the buffer size for very large MTL files if necessary
	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	var currentMaterial *Material
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			// Skip empty lines and comments
			continue
		}
		// Find the first space to determine the directive
		firstSpace := strings.IndexByte(line, ' ')
		if firstSpace == -1 {
			// Line has only one token, possibly invalid
			continue // Or handle single-token directives if any
		}
		directive := line[:firstSpace]
		rest := line[firstSpace+1:]

		switch directive {
		case "newmtl":
			if rest == "" {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "newmtl directive error: material name is missing",
				}
			}
			name := rest
			mat := &Material{Name: name}
			p.obj.Materials[name] = mat
			currentMaterial = mat
		case "Kd":
			if currentMaterial == nil {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Kd directive error: defined before any newmtl",
				}
			}
			fields := splitFields(rest, 3)
			if len(fields) < 3 {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Kd directive error: expected 3 components",
				}
			}
			for i := 0; i < 3; i++ {
				val, err := strconv.ParseFloat(fields[i], 64)
				if err != nil {
					return &ParseError{
						Filename: mtlFilename,
						Line:     lineNumber,
						LineText: line,
						Msg:      fmt.Sprintf("invalid Kd value: %v", err),
					}
				}
				currentMaterial.Diffuse[i] = val
			}
		case "Ka":
			if currentMaterial == nil {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Ka directive error: defined before any newmtl",
				}
			}
			fields := splitFields(rest, 3)
			if len(fields) < 3 {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Ka directive error: expected 3 components",
				}
			}
			for i := 0; i < 3; i++ {
				val, err := strconv.ParseFloat(fields[i], 64)
				if err != nil {
					return &ParseError{
						Filename: mtlFilename,
						Line:     lineNumber,
						LineText: line,
						Msg:      fmt.Sprintf("invalid Ka value: %v", err),
					}
				}
				currentMaterial.Ambient[i] = val
			}
		case "Ks":
			if currentMaterial == nil {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Ks directive error: defined before any newmtl",
				}
			}
			fields := splitFields(rest, 3)
			if len(fields) < 3 {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Ks directive error: expected 3 components",
				}
			}
			for i := 0; i < 3; i++ {
				val, err := strconv.ParseFloat(fields[i], 64)
				if err != nil {
					return &ParseError{
						Filename: mtlFilename,
						Line:     lineNumber,
						LineText: line,
						Msg:      fmt.Sprintf("invalid Ks value: %v", err),
					}
				}
				currentMaterial.Specular[i] = val
			}
		case "Ns":
			if currentMaterial == nil {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Ns directive error: defined before any newmtl",
				}
			}
			fields := splitFields(rest, 1)
			if len(fields) < 1 {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "Ns directive error: expected a single value",
				}
			}
			val, err := strconv.ParseFloat(fields[0], 64)
			if err != nil {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      fmt.Sprintf("invalid Ns value: %v", err),
				}
			}
			currentMaterial.Shininess = val
		case "map_Kd":
			if currentMaterial == nil {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "map_Kd directive error: defined before any newmtl",
				}
			}
			if rest == "" {
				return &ParseError{
					Filename: mtlFilename,
					Line:     lineNumber,
					LineText: line,
					Msg:      "map_Kd directive error: expected a filename",
				}
			}
			texture := rest
			currentMaterial.Texture = texture
		default:
			// Ignore other material properties or unsupported directives.
		}
	}
	if err := scanner.Err(); err != nil {
		return &ParseError{
			Filename: mtlFilename,
			Line:     lineNumber,
			LineText: "",
			Msg:      fmt.Sprintf("error reading MTL file: %v", err),
		}
	}
	return nil
}

// splitFields splits a string into fields, expecting exactly 'limit' fields.
// If limit is negative, it returns all fields.
func splitFields(s string, limit int) []string {
	if limit > 0 {
		return strings.SplitN(s, " ", limit)
	}
	return strings.Fields(s)
}

// newError creates a new ParseError with the current parser state and a formatted message.
func (p *parser) newError(format string, args ...interface{}) error {
	return &ParseError{
		Filename: p.filename,
		Line:     p.lineNumber,
		LineText: p.lineText,
		Msg:      fmt.Sprintf(format, args...),
	}
}

func EncodeOBJ(w io.Writer, obj *Object) error {
	if err := obj.WriteOBJ(w); err != nil {
		return fmt.Errorf("failed to write OBJ data: %v", err)
	}
	if err := obj.WriteMTL(w); err != nil {
		return fmt.Errorf("failed to write MTL data: %v", err)
	}
	return nil
}

// WriteOBJ serializes the Object to an OBJ file.
// It writes vertex data, face definitions, and material references to the provided writer.
func (obj *Object) WriteOBJ(w io.Writer) error {
	writer := bufio.NewWriter(w)
	if len(obj.Materials) > 0 {
		if _, err := fmt.Fprintf(writer, "mtllib materials.mtl\n"); err != nil {
			return fmt.Errorf("failed to write mtllib: %v", err)
		}
	}
	// Preallocate if possible.
	vertices := obj.Vertices
	texCoords := obj.TexCoords
	normals := obj.Normals
	faces := obj.Faces

	for _, v := range vertices {
		if _, err := fmt.Fprintf(writer, "v %f %f %f\n", v.X, v.Y, v.Z); err != nil {
			return fmt.Errorf("failed to write vertex: %v", err)
		}
	}
	for _, vt := range texCoords {
		if _, err := fmt.Fprintf(writer, "vt %f %f\n", vt.U, vt.V); err != nil {
			return fmt.Errorf("failed to write texture coordinate: %v", err)
		}
	}
	for _, vn := range normals {
		if _, err := fmt.Fprintf(writer, "vn %f %f %f\n", vn.X, vn.Y, vn.Z); err != nil {
			return fmt.Errorf("failed to write normal: %v", err)
		}
	}
	// Keep track of the current material to write 'usemtl' only when it changes.
	var currentMaterial string
	for _, face := range faces {
		// Write 'usemtl' if the material changes.
		if face.Material != currentMaterial {
			if face.Material != "" {
				if _, err := fmt.Fprintf(writer, "usemtl %s\n", face.Material); err != nil {
					return fmt.Errorf("failed to write usemtl: %v", err)
				}
			}
			currentMaterial = face.Material
		}

		// Write face indices.
		var faceLine strings.Builder
		faceLine.Grow(64) // Estimate initial capacity
		faceLine.WriteString("f")
		for _, idx := range face.Indices {
			faceLine.WriteString(" ")
			faceLine.WriteString(formatIndex(idx))
		}
		faceLine.WriteString("\n")
		if _, err := writer.WriteString(faceLine.String()); err != nil {
			return fmt.Errorf("failed to write face: %v", err)
		}
	}

	return writer.Flush()
}

// WriteMTL serializes the Materials map to an MTL file.
// It writes material definitions to the provided writer.
func (obj *Object) WriteMTL(w io.Writer) error {
	if len(obj.Materials) == 0 {
		return nil // No materials to write.
	}
	writer := bufio.NewWriter(w)

	for _, mat := range obj.Materials {
		if mat == nil {
			continue // Skip nil materials
		}
		if _, err := fmt.Fprintf(writer, "newmtl %s\n", mat.Name); err != nil {
			return fmt.Errorf("failed to write newmtl for material '%s': %v", mat.Name, err)
		}
		if mat.Diffuse != [3]float64{} {
			if _, err := fmt.Fprintf(writer, "Kd %f %f %f\n", mat.Diffuse[0], mat.Diffuse[1], mat.Diffuse[2]); err != nil {
				return fmt.Errorf("failed to write Kd for material '%s': %v", mat.Name, err)
			}
		}
		if mat.Ambient != [3]float64{} {
			if _, err := fmt.Fprintf(writer, "Ka %f %f %f\n", mat.Ambient[0], mat.Ambient[1], mat.Ambient[2]); err != nil {
				return fmt.Errorf("failed to write Ka for material '%s': %v", mat.Name, err)
			}
		}
		if mat.Specular != [3]float64{} {
			if _, err := fmt.Fprintf(writer, "Ks %f %f %f\n", mat.Specular[0], mat.Specular[1], mat.Specular[2]); err != nil {
				return fmt.Errorf("failed to write Ks for material '%s': %v", mat.Name, err)
			}
		}
		if mat.Shininess != 0 {
			if _, err := fmt.Fprintf(writer, "Ns %f\n", mat.Shininess); err != nil {
				return fmt.Errorf("failed to write Ns for material '%s': %v", mat.Name, err)
			}
		}
		if mat.Texture != "" {
			if _, err := fmt.Fprintf(writer, "map_Kd %s\n", mat.Texture); err != nil {
				return fmt.Errorf("failed to write map_Kd for material '%s': %v", mat.Name, err)
			}
		}
		if _, err := writer.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline after material '%s': %v", mat.Name, err)
		}
	}

	return writer.Flush()
}

// formatIndex formats an Index into the OBJ face index format.
func formatIndex(idx Index) string {
	// OBJ indices are 1-based.
	var parts []string
	v := strconv.Itoa(idx.Vertex)
	parts = append(parts, v)
	if idx.TexCoord != 0 || idx.Normal != 0 {
		vt := ""
		if idx.TexCoord != 0 {
			vt = strconv.Itoa(idx.TexCoord)
		}
		vn := ""
		if idx.Normal != 0 {
			vn = strconv.Itoa(idx.Normal)
		}
		parts = append(parts, vt)
		parts = append(parts, vn)
		return strings.Join(parts, "/")
	}
	return v
}
