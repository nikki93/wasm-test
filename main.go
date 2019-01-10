package main

import (
	"syscall/js"
	"wasm-test/webgl"
)

func main() {
	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "main_canvas")
	devicePixelRatio := js.Global().Get("devicePixelRatio").Float()
	width := int(devicePixelRatio * doc.Get("body").Get("clientWidth").Float())
	height := int(devicePixelRatio * doc.Get("body").Get("clientHeight").Float())
	canvas.Set("width", width)
	canvas.Set("height", height)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false
	gl, _ := webgl.NewContext(canvas, attrs)

	//// VERTEX BUFFER ////
	var verticesNative = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	var vertices = js.TypedArrayOf(verticesNative)
	// Create buffer
	vertexBuffer := gl.CreateBuffer()
	// Bind to buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)

	// Pass data to buffer
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	// Unbind buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, js.Null())

	//// INDEX BUFFER ////
	var indicesNative = []uint32{
		2, 1, 0,
	}
	var indices = js.TypedArrayOf(indicesNative)

	// Create buffer
	indexBuffer := gl.CreateBuffer()

	// Bind to buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)

	// Pass data to buffer
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

	// Unbind buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, js.Null())

	//// Shaders ////

	// Vertex shader source code
	vertCode := `
	attribute vec3 coordinates;
		
	void main(void) {
		gl_Position = vec4(coordinates, 1.0);
	}`

	// Create a vertex shader object
	vertShader := gl.CreateShader(gl.VERTEX_SHADER)

	// Attach vertex shader source code
	gl.ShaderSource(vertShader, vertCode)

	// Compile the vertex shader
	gl.CompileShader(vertShader)

	//fragment shader source code
	fragCode := `
	void main(void) {
		gl_FragColor = vec4(1.0, 0.0, 1.0, 1.0);
	}`

	// Create fragment shader object
	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)

	// Attach fragment shader source code
	gl.ShaderSource(fragShader, fragCode)

	// Compile the fragmentt shader
	gl.CompileShader(fragShader)

	// Create a shader program object to store
	// the combined shader program
	shaderProgram := gl.Call("createProgram")

	// Attach a vertex shader
	gl.AttachShader(shaderProgram, vertShader)

	// Attach a fragment shader
	gl.AttachShader(shaderProgram, fragShader)

	// Link both the programs
	gl.LinkProgram(shaderProgram)

	// Use the combined shader program object
	gl.UseProgram(shaderProgram)

	//// Associating shaders to buffer objects ////

	// Bind vertex buffer object
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)

	// Bind index buffer object
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)

	// Get the attribute location
	coord := gl.GetAttribLocation(shaderProgram, "coordinates")

	// Point an attribute to the currently bound VBO
	gl.VertexAttribPointer(coord, 3, gl.FLOAT, false, 0, 0)

	// Enable the attribute
	gl.EnableVertexAttribArray(coord)

	//// Drawing the triangle ////

	// Clear the canvas
	gl.ClearColor(0.5, 0.5, 0.5, 0.9)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Enable the depth test
	gl.Enable(gl.DEPTH_TEST)

	// Set the view port
	gl.Viewport(0, 0, width, height)

	// Draw the triangle
	gl.DrawElements(gl.TRIANGLES, len(indicesNative), gl.UNSIGNED_SHORT, 0)
}
