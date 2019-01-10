package main

import (
	"syscall/js"
	"wasm-test/webgl"
)

func main() {
	// Canvas
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "main_canvas")
	devicePixelRatio := js.Global().Get("devicePixelRatio").Float()
	width := int(devicePixelRatio * doc.Get("body").Get("clientWidth").Float())
	height := int(devicePixelRatio * doc.Get("body").Get("clientHeight").Float())
	canvas.Set("width", width)
	canvas.Set("height", height)

	// GL
	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false
	gl, _ := webgl.NewContext(canvas, attrs)

	// Vertex buffer
	var verticesNative = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	var vertices = js.TypedArrayOf(verticesNative)
	vertexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, js.Null())

	// Index buffer
	var indicesNative = []uint32{
		2, 1, 0,
	}
	var indices = js.TypedArrayOf(indicesNative)
	indexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, js.Null())

	// Vertex shader
	vert := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vert, `
		attribute vec3 coordinates;
			
		void main(void) {
			gl_Position = vec4(coordinates, 1.0);
		}
	`)
	gl.CompileShader(vert)

	// Fragment shader
	frag := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(frag, `
		void main(void) {
			gl_FragColor = vec4(1.0, 0.0, 1.0, 1.0);
		}
	`)
	gl.CompileShader(frag)

	// Program
	program := gl.CreateProgram()
	gl.AttachShader(program, vert)
	gl.AttachShader(program, frag)
	gl.LinkProgram(program)
	gl.UseProgram(program)

	// Attributes
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	coord := gl.GetAttribLocation(program, "coordinates")
	gl.VertexAttribPointer(coord, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(coord)

	// Draw
	gl.ClearColor(0.5, 0.5, 0.5, 0.9)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.Viewport(0, 0, width, height)
	gl.DrawElements(gl.TRIANGLES, len(indicesNative), gl.UNSIGNED_SHORT, 0)
}
