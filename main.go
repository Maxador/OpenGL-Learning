package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"runtime"
	"strings"
)

const width = 800
const height = 600

func init() {
	runtime.LockOSThread()
}

func main() {
	defer glfw.Terminate()

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	// var vertices = []float32{
	// 	0.5, 0.5, 0.0, // Top right
	// 	0.5, -0.5, 0.0, // Bottom right
	// 	-0.5, -0.5, 0.0, // Bottom left
	// 	-0.5, 0.5, 0.0, // Top left
	// }

	// var vertices = []float32{
	// 	// First triangle
	// 	0.0, -0.5, 0.0,
	// 	0.45, 0.5, 0.0,
	// 	0.9, -0.5, 0.0,
	//
	// 	// Second triangle
	// 	0.0, -0.5, 0.0,
	// 	-0.9, -0.5, 0.0,
	// 	-0.45, 0.5, 0.0,
	// }

	var verticesLeft = []float32{
		0.0, -0.5, 0.0,
		-0.9, -0.5, 0.0,
		-0.45, 0.5, 0.0,
	}
	var verticesRight = []float32{
		0.0, -0.5, 0.0,
		0.45, 0.5, 0.0,
		0.9, -0.5, 0.0,
	}

	// var indices = []uint32{
	// 	0, 1, 3, // First triangle
	// 	1, 2, 3, // Second triangle
	// }

	// var vao, vbo, ebo uint32
	var vao, vbo [2]uint32
	gl.GenVertexArrays(1, &vao[0])
	gl.GenVertexArrays(1, &vao[1])
	gl.GenBuffers(1, &vbo[0])
	gl.GenBuffers(1, &vbo[1])
	// gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(verticesRight)*4, gl.Ptr(verticesRight), gl.STATIC_DRAW)

	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	// Triangle 2

	gl.BindVertexArray(vao[1])
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(verticesLeft)*4, gl.Ptr(verticesLeft), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	program, err := createProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(vao[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(vao[1])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("Failed to compile %v : %v", source, log)
	}
	return shader, nil
}

func createProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("Failed to link the program: %v", log)
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return program, nil
}

var vertexShaderSource = `
#version 410 core

layout (location = 0) in vec3 position;

void main() {
	gl_Position = vec4(position.x, position.y, position.z, 1.0);
}
` + "\x00"

var fragmentShaderSource = `
#version 410 core

out vec4 color;

void main() {
	color = vec4(1.0, 0.5, 0.2, 1.0);
}` + "\x00"
