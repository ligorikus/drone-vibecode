//go:build !no_vis

package visualization

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"drone/internal/config"
	"drone/internal/models"
	"drone/internal/services"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	// Shaders for rendering drones
	vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 position;
uniform mat4 MVP;
void main()
{
	gl_Position = MVP * vec4(position, 1.0);
}
` + "\x00"

	fragmentShaderSource = `
#version 330 core
out vec4 color;
uniform vec3 droneColor;
void main()
{
	color = vec4(droneColor, 1.0);
}
` + "\x00"

	// Shaders for ground plane with Perlin noise
	groundVertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 position;
layout (location = 1) in vec2 texCoord;
out vec2 fragTexCoord;
uniform mat4 MVP;
void main()
{
	gl_Position = MVP * vec4(position, 1.0);
	fragTexCoord = texCoord;
}
` + "\x00"

	groundFragmentShaderSource = `
#version 330 core
out vec4 color;
in vec2 fragTexCoord;

// Simple hash function
float hash(vec2 p) {
	return fract(sin(dot(p, vec2(127.1, 311.7))) * 43758.5453123);
}

// 2D Perlin noise
float noise(vec2 p) {
	vec2 i = floor(p);
	vec2 f = fract(p);
	
	// Smoothstep for interpolation
	vec2 u = f * f * (3.0 - 2.0 * f);
	
	// Get corner values
	float a = hash(i);
	float b = hash(i + vec2(1.0, 0.0));
	float c = hash(i + vec2(0.0, 1.0));
	float d = hash(i + vec2(1.0, 1.0));
	
	// Mix
	return mix(mix(a, b, u.x), mix(c, d, u.x), u.y);
}

// Fractal Brownian Motion
float fbm(vec2 p) {
	float value = 0.0;
	float amplitude = 0.5;
	float frequency = 1.0;
	
	for (int i = 0; i < 5; i++) {
		value += amplitude * noise(p * frequency);
		frequency *= 2.0;
		amplitude *= 0.5;
	}
	return value;
}

void main()
{
	// Sample noise at texture coordinates with scaling
	vec2 uv = fragTexCoord;
	float n = fbm(uv * 8.0);
	
	// Create colors for visibility - NO cyan/blue, only warm colors
	vec3 darkColor = vec3(0.08, 0.08, 0.1);       // Dark background
	vec3 accentColor = vec3(0.8, 0.6, 0.0);        // Orange-gold accent (NOT cyan)
	vec3 highlightColor = vec3(1.0, 0.9, 0.3);     // Bright yellow
	
	// Layer the colors based on noise value
	vec3 groundColor;
	if (n > 0.6) {
		groundColor = mix(accentColor, highlightColor, (n - 0.6) / 0.4);
	} else if (n > 0.35) {
		groundColor = mix(darkColor, accentColor, (n - 0.35) / 0.25);
	} else {
		groundColor = darkColor;
	}
	
	// Add prominent grid lines for better motion perception
	float gridSize = 10.0;
	vec2 gridUV = uv * gridSize;
	float gridX = abs(fract(gridUV.x) - 0.5);
	float gridY = abs(fract(gridUV.y) - 0.5);
	float grid = min(gridX, gridY);
	float gridMask = smoothstep(0.03, 0.0, grid);
	groundColor = mix(groundColor, vec3(0.7, 0.7, 0.7), gridMask * 0.5);
	
	// Add axis indicators (red for +X, green for +Z) - NOT blue
	float axisThickness = 0.05;
	float xAxis = smoothstep(axisThickness, 0.0, abs(uv.y - 0.5));
	float zAxis = smoothstep(axisThickness, 0.0, abs(uv.x - 0.5));
	if (uv.x > 0.5) {
		groundColor = mix(groundColor, vec3(1.0, 0.2, 0.2), xAxis * 0.7);
	}
	if (uv.y > 0.5) {
		groundColor = mix(groundColor, vec3(0.2, 0.8, 0.2), zAxis * 0.7);
	}
	
	color = vec4(groundColor, 1.0);
}
` + "\x00"
)

// OpenGLVisualizer реализует Visualizer с использованием OpenGL
type OpenGLVisualizer struct {
	window       *glfw.Window
	mainDrone    *models.Drone
	config       *config.Config
	children     []*models.ChildDrone
	program      uint32
	vao          uint32
	vbo          uint32
	mvpUniform   int32
	colorUniform int32

	// Ground plane resources
	groundProgram    uint32
	groundVAO        uint32
	groundVBO        uint32
	groundTexCoordLoc int32

	// Camera control
	cameraDistance float32
	cameraAngleX   float32
	cameraAngleY   float32
	lastMouseX     float64
	lastMouseY     float64
	mousePressed   bool

	// Simulation service for drone control
	simulation services.SimulationProvider

	// Keyboard input state
	keys map[glfw.Key]bool
}

// NewOpenGLVisualizer создаёт новый OpenGL визуализатор
func NewOpenGLVisualizer(mainDrone *models.Drone, cfg *config.Config) *OpenGLVisualizer {
	return &OpenGLVisualizer{
		mainDrone:      mainDrone,
		config:         cfg,
		cameraDistance: 40.0,  // Увеличил дистанцию
		cameraAngleX:   20.0,  // Более горизонтальный угол (было 45°)
		cameraAngleY:   45.0,
		keys:           make(map[glfw.Key]bool),
	}
}

// SetSimulation устанавливает сервис симуляции для управления дронами
func (ov *OpenGLVisualizer) SetSimulation(sim services.SimulationProvider) {
	ov.simulation = sim
}

// Init инициализирует OpenGL контекст и ресурсы
func (ov *OpenGLVisualizer) Init() error {
	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize GLFW: %v", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create window
	window, err := glfw.CreateWindow(ov.config.WindowWidth, ov.config.WindowHeight, "Drone Simulation Visualization", nil, nil)
	if err != nil {
		glfw.Terminate()
		return fmt.Errorf("failed to create GLFW window: %v", err)
	}
	ov.window = window

	// Make context current
	window.MakeContextCurrent()

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL: %v", err)
	}

	// Enable depth testing
	gl.Enable(gl.DEPTH_TEST)

	// Setup mouse callbacks
	ov.setupMouseCallbacks()

	// Compile shaders
	vertexShader := ov.compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := ov.compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	// Create shader program
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	// Check for linking errors
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return fmt.Errorf("failed to link shader program: %s", log[:logLength])
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	ov.program = program
	ov.mvpUniform = gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	ov.colorUniform = gl.GetUniformLocation(program, gl.Str("droneColor\x00"))

	// Create VAO and VBO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	ov.vao = vao

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	ov.vbo = vbo

	// Define cube vertices for drone representation (12 triangles for 6 faces)
	// Smaller size for better visibility
	droneSize := float32(0.3)
	vertices := []float32{
		// Front face
		-droneSize, -droneSize, droneSize,
		droneSize, -droneSize, droneSize,
		droneSize, droneSize, droneSize,
		-droneSize, -droneSize, droneSize,
		droneSize, droneSize, droneSize,
		-droneSize, droneSize, droneSize,
		// Back face
		-droneSize, -droneSize, -droneSize,
		-droneSize, droneSize, -droneSize,
		droneSize, droneSize, -droneSize,
		-droneSize, -droneSize, -droneSize,
		droneSize, droneSize, -droneSize,
		droneSize, -droneSize, -droneSize,
		// Top face
		-droneSize, droneSize, -droneSize,
		-droneSize, droneSize, droneSize,
		droneSize, droneSize, droneSize,
		-droneSize, droneSize, -droneSize,
		droneSize, droneSize, droneSize,
		droneSize, droneSize, -droneSize,
		// Bottom face
		-droneSize, -droneSize, -droneSize,
		droneSize, -droneSize, -droneSize,
		droneSize, -droneSize, droneSize,
		-droneSize, -droneSize, -droneSize,
		droneSize, -droneSize, droneSize,
		-droneSize, -droneSize, droneSize,
		// Right face
		droneSize, -droneSize, -droneSize,
		droneSize, droneSize, -droneSize,
		droneSize, droneSize, droneSize,
		droneSize, -droneSize, -droneSize,
		droneSize, droneSize, droneSize,
		droneSize, -droneSize, droneSize,
		// Left face
		-droneSize, -droneSize, -droneSize,
		-droneSize, -droneSize, droneSize,
		-droneSize, droneSize, droneSize,
		-droneSize, -droneSize, -droneSize,
		-droneSize, droneSize, droneSize,
		-droneSize, droneSize, -droneSize,
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Initialize ground plane
	if err := ov.initGround(); err != nil {
		return fmt.Errorf("failed to initialize ground: %v", err)
	}

	return nil
}

func (ov *OpenGLVisualizer) compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		panic(fmt.Sprintf("failed to compile %s shader: %s",
			map[uint32]string{gl.VERTEX_SHADER: "vertex", gl.FRAGMENT_SHADER: "fragment"}[shaderType],
			log[:logLength]))
	}

	return shader
}

// initGround инициализирует плоскость земли с текстурой Перлина
func (ov *OpenGLVisualizer) initGround() error {
	// Compile ground shaders
	groundVertexShader := ov.compileShader(groundVertexShaderSource, gl.VERTEX_SHADER)
	groundFragmentShader := ov.compileShader(groundFragmentShaderSource, gl.FRAGMENT_SHADER)

	// Create ground shader program
	groundProgram := gl.CreateProgram()
	gl.AttachShader(groundProgram, groundVertexShader)
	gl.AttachShader(groundProgram, groundFragmentShader)
	gl.LinkProgram(groundProgram)

	// Check for linking errors
	var status int32
	gl.GetProgramiv(groundProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(groundProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(groundProgram, logLength, nil, gl.Str(log))
		return fmt.Errorf("failed to link ground shader program: %s", log[:logLength])
	}

	gl.DeleteShader(groundVertexShader)
	gl.DeleteShader(groundFragmentShader)

	ov.groundProgram = groundProgram
	ov.groundTexCoordLoc = gl.GetUniformLocation(groundProgram, gl.Str("MVP\x00"))

	// Create ground VAO
	var groundVAO uint32
	gl.GenVertexArrays(1, &groundVAO)
	gl.BindVertexArray(groundVAO)
	ov.groundVAO = groundVAO

	// Create ground VBO with position and texCoord data
	// Large ground plane centered at origin (XZ plane, y=0)
	groundSize := float32(500.0)

	// Vertices: position (x, y, z) + texCoord (u, v)
	groundVertices := []float32{
		// Front-left
		-groundSize, 0.0, -groundSize, 0.0, 0.0,
		// Front-right
		groundSize, 0.0, -groundSize, 1.0, 0.0,
		// Back-right
		groundSize, 0.0, groundSize, 1.0, 1.0,
		// Front-left
		-groundSize, 0.0, -groundSize, 0.0, 0.0,
		// Back-right
		groundSize, 0.0, groundSize, 1.0, 1.0,
		// Back-left
		-groundSize, 0.0, groundSize, 0.0, 1.0,
	}

	var groundVBO uint32
	gl.GenBuffers(1, &groundVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, groundVBO)
	ov.groundVBO = groundVBO

	gl.BufferData(gl.ARRAY_BUFFER, len(groundVertices)*4, gl.Ptr(groundVertices), gl.STATIC_DRAW)

	// Position attribute (location = 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// TexCoord attribute (location = 1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	return nil
}

// RenderLoop запускает цикл рендеринга
func (ov *OpenGLVisualizer) RenderLoop(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 16) // ~60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if ov.window.ShouldClose() {
				return nil
			}

			if err := ov.render(); err != nil {
				return err
			}
		}
	}
}

func (ov *OpenGLVisualizer) render() error {
	// Process keyboard input and update main drone
	ov.processInput(0.016) // ~60 FPS

	// Clear screen with dark background matching ground
	gl.ClearColor(0.05, 0.05, 0.08, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Update drone positions and render
	ov.renderDrones()

	// Swap buffers
	ov.window.SwapBuffers()

	// Poll events
	glfw.PollEvents()

	return nil
}

func (ov *OpenGLVisualizer) renderDrones() {
	// Get camera position from mouse control
	// Camera follows the main drone
	mainPos := ov.mainDrone.GetPosition()
	target := mgl32.Vec3{float32(mainPos.X), float32(mainPos.Y), float32(mainPos.Z)}

	cameraPos := ov.getCameraPosition(target)

	// Set up camera/view matrix
	view := mgl32.LookAt(
		cameraPos[0], cameraPos[1], cameraPos[2],
		target[0], target[1], target[2],
		0, 1, 0, // Up vector
	)

	// Set up projection matrix
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(ov.config.WindowWidth)/float32(ov.config.WindowHeight), 0.1, 1000.0)

	// Render ground FIRST (with Perlin noise texture) - static ground at Y=0
	ov.renderGround(projection, view)

	// Bind drone VAO
	gl.BindVertexArray(ov.vao)

	// Use drone shader program
	gl.UseProgram(ov.program)

	// Main drone (red)
	model := mgl32.Translate3D(float32(mainPos.X), float32(mainPos.Y), float32(mainPos.Z))
	mvp := projection.Mul4(view).Mul4(model)

	gl.UniformMatrix4fv(ov.mvpUniform, 1, false, &mvp[0])
	gl.Uniform3f(ov.colorUniform, 1.0, 0.0, 0.0) // Red color for main drone
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	// Child drones (blue)
	children := ov.mainDrone.GetChildren()
	for _, child := range children {
		childPos := child.GetPosition()
		model = mgl32.Translate3D(float32(childPos.X), float32(childPos.Y), float32(childPos.Z))
		mvp = projection.Mul4(view).Mul4(model)

		gl.UniformMatrix4fv(ov.mvpUniform, 1, false, &mvp[0])
		gl.Uniform3f(ov.colorUniform, 0.0, 0.5, 1.0) // Blue color for child drones
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
}

// renderGround рендерит плоскость земли с текстурой Перлина
func (ov *OpenGLVisualizer) renderGround(projection, view mgl32.Mat4) {
	// Enable depth test for ground
	gl.Enable(gl.DEPTH_TEST)
	
	gl.UseProgram(ov.groundProgram)

	// Ground model matrix - static at origin (y=0), no offset
	// This keeps the ground fixed so drone movement is visible
	model := mgl32.Mat4{}
	model[0] = 1.0
	model[5] = 1.0
	model[10] = 1.0
	model[15] = 1.0
	
	mvp := projection.Mul4(view).Mul4(model)

	gl.UniformMatrix4fv(ov.groundTexCoordLoc, 1, false, &mvp[0])

	// Bind ground VAO and draw
	gl.BindVertexArray(ov.groundVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

// Close закрывает ресурсы визуализатора
func (ov *OpenGLVisualizer) Close() error {
	if ov.window != nil {
		ov.window.Destroy()
	}
	
	// Delete ground resources
	if ov.groundVAO != 0 {
		gl.DeleteVertexArrays(1, &ov.groundVAO)
	}
	if ov.groundVBO != 0 {
		gl.DeleteBuffers(1, &ov.groundVBO)
	}
	if ov.groundProgram != 0 {
		gl.DeleteProgram(ov.groundProgram)
	}
	
	glfw.Terminate()
	return nil
}

// GetMainDrone возвращает главный дрон
func (ov *OpenGLVisualizer) GetMainDrone() *models.Drone {
	return ov.mainDrone
}

func (ov *OpenGLVisualizer) setupMouseCallbacks() {
	// Keyboard callback for drone control
	ov.window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			ov.keys[key] = true
		} else if action == glfw.Release {
			ov.keys[key] = false
		}
	})

	// Mouse button callback
	ov.window.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				ov.mousePressed = true
				ov.lastMouseX, ov.lastMouseY = window.GetCursorPos()
			} else if action == glfw.Release {
				ov.mousePressed = false
			}
		}
	})

	// Cursor position callback for rotation
	ov.window.SetCursorPosCallback(func(window *glfw.Window, xpos float64, ypos float64) {
		if ov.mousePressed {
			dx := xpos - ov.lastMouseX
			dy := ypos - ov.lastMouseY

			ov.cameraAngleY += float32(dx) * 0.5
			ov.cameraAngleX += float32(dy) * 0.5

			// Limit vertical angle to avoid flipping
			if ov.cameraAngleX > 89.0 {
				ov.cameraAngleX = 89.0
			}
			if ov.cameraAngleX < -89.0 {
				ov.cameraAngleX = -89.0
			}

			ov.lastMouseX = xpos
			ov.lastMouseY = ypos
		}
	})

	// Scroll callback for zoom
	ov.window.SetScrollCallback(func(window *glfw.Window, xoff float64, yoff float64) {
		ov.cameraDistance -= float32(yoff) * 2.0

		// Limit zoom range
		if ov.cameraDistance < 5.0 {
			ov.cameraDistance = 5.0
		}
		if ov.cameraDistance > 100.0 {
			ov.cameraDistance = 100.0
		}
	})
}

func (ov *OpenGLVisualizer) getCameraPosition(target mgl32.Vec3) mgl32.Vec3 {
	// Convert spherical coordinates to Cartesian
	radiansX := float64(ov.cameraAngleX) * math.Pi / 180.0
	radiansY := float64(ov.cameraAngleY) * math.Pi / 180.0
	distance := float64(ov.cameraDistance)

	x := float32(distance * math.Cos(radiansX) * math.Cos(radiansY))
	y := float32(distance * math.Sin(radiansX))
	z := float32(distance * math.Cos(radiansX) * math.Sin(radiansY))

	// Camera position relative to target
	return mgl32.Vec3{
		target[0] + x,
		target[1] + y,
		target[2] + z,
	}
}

// processInput обрабатывает ввод с клавиатуры и обновляет главного дрона
func (ov *OpenGLVisualizer) processInput(deltaTime float64) {
	if ov.simulation == nil {
		return
	}

	input := services.InputState{
		Forward:  ov.keys[glfw.KeyW],
		Backward: ov.keys[glfw.KeyS],
		Left:     ov.keys[glfw.KeyA],
		Right:    ov.keys[glfw.KeyD],
		Up:       ov.keys[glfw.KeySpace],
		Down:     ov.keys[glfw.KeyLeftShift] || ov.keys[glfw.KeyRightShift],
	}

	ov.simulation.SetInput(input)
	ov.simulation.UpdateMainDrone(deltaTime)
}
