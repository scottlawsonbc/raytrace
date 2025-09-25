# raytrace — A Physically Based Path Tracer in Go (standard library only)

I wrote this renderer to learn graphics the hard way:
by building a physically based path tracer in pure Go,
with no dependencies beyond the standard library.
If you are curious about rendering, want a readable codebase,
and like to tinker, this project is for you.

The core lives in the `phys` package.
Scenes are small Go programs in `scene/` you can run directly.
I kept a dev journal over three years; the story shows in the code.

---

## Highlights

- **Physically based** path tracer with emissive lights, textured meshes, and a clean `phys` core.
- **BVH with binned SAH** splitter for fast intersection; big speedups over naive splits.
- **OBJ/MTL + PNG textures** loader for real-world assets and 3D scans.
- **Deterministic, dependency-free Go**: only the standard library.

---

## Gallery

_You will find lots of visuals in `scene/`._
Add your images and gifs to the repo so GitHub renders them inline.

- `scene/gem/` – caustics & refraction
  `./scene/gem/out/out_20241106_000741.gif`
- `scene/scan/` – textured 3D scans, camera orbit
  `./scene/scan/out/out_20241119_161650.gif`

---

## Quick start

```bash
# Go 1.21+ recommended
git clone <your-repo-url>
cd raytrace

# Run a scene
cd scene/gem
go run .

# Or run by path
go run ./scene/scan


# raytrace

Simple raytracer written in Go. Started working on this in fall 2022.
Continued more work in fall 2024.

Loosely based on the book "Ray Tracing in One Weekend" by Peter Shirley.

The main raytracer code is in the `phys` package.
Other modules can call this to setup a scene and render it.

<img src="./scene/gem/out/out_20241106_000741.gif" height="400px">

One of the nice things I like about raytracer is that it only uses the standard
library in Go, and that means I can fully understand every line of code in the
codebase. Although it does such crazy operations with light, and linear algebra,
it has been a very rewarding experience to write each piece of code and see it
all come together.


## Usage

Run one of the scenes in the `slam/code/photon/raytrace/scene` folder.
For example, `go run .` in the directory or by calling `go run` and pointing
it directly to the main.go or other files.


## Coordinates

Here are some reference notes about the different coordinate systems used in
this or other raytracers.


World coordinates means phys.Distance and is the main coordinate system used in
the raytracer. It is the space where the scene is defined.

Camera coordinates are the local coordinates attached to the camera with axes
aligned to match the screen orientation.

Screen coordinates are the 2D coordinates of the screen where the image is
rendered. The screen is defined by the camera and the distance to the screen.

Coordinate systems:
- World [X Y Z] The space where the scene is defined.
- Camera [S, T]: The space where the camera is defined.
- Screen []: The space where the screen is defined.
- Texture space: The space where the texture is defined.
- Object space: The space where the object local is defined.

An alternative unambiguous set of coordinates can be done with two letters.
Virtually every other renderer uses one letter, but I think two letters is
suitable because there are a few coordinate systems where it would be nice
to re-use the letters as it is more clear, x, y, z for example in world space,
but also x, y, z in object space, and x, y in screen space.
Likewise, for camera ray generation and for texture mapping we can use u, v,
prefixed with the first letter of the coordinate system.

So with our second letter, we can have:

- sx, sy: Screen
- vx, vy, vz: View
- wx, wy, wz: World
- lx, ly, lz: Local
- tu, tv: Texture
- cu, cv: Camera



# 2024-11-10 thoughts about next steps

I have just finished a lot of work overhauling the raytracer completely, it is
now a path tracer with a lot more functionality.

I want to improve my direct lighting support and support realistic materials
by implementing a BSDF interface for my different materials.

I still also haven't investigated importance sampling and integrators yet.

to implement:
- BxDF interface for materials
- Directional light source (for narrow beams, etc)
- integrator interface
- img diffing tools and tests
- go wasm for frontend as well, to leverage the types in phys


# 2024-11-11 suggest debug shaders


Depth Shader: Visualizes depth from the camera to each pixel.
UV Coordinate Shader: Maps UV coordinates to colors for texture mapping verification.
Albedo Shader: Displays the base color of materials without lighting.
World Position Shader: Colors pixels based on their world space positions.
Emission Shader: Highlights emissive materials and light sources.
Wireframe Shader: Renders the mesh's wireframe to inspect geometry.
Tangent Space Shader: Visualizes tangent, bitangent, and normal vectors.
Emission Intensity Shader: Shows the intensity of emissive materials.
Ambient Occlusion Shader: Displays ambient occlusion factors as grayscale.
Albedo and Normal Overlay Shader: Combines albedo and normal visualizations.


I can realistically do the depth shader, uv coordinate shader, albedo shader,
world position shader, emission shader, wireframe shader, and tangent space.
It seems these shaders might need to do interface checks since not all of the
materials have a base color or albedo.

type ShaderDepth struct {}
type ShaderUV struct {}
type ShaderTangentSpace struct {}


# 2024-11-12 things to improve next

1. dielectric doesn't seem to work at all when roughness = 0. I'm not sure why.
2. texture coordinates may be inverted, I'm not sure, I don't know enough yet.


# 2024-11-12 how does geometric algebra fit in?

I've been reminding myself that even though I don't quite grasp it yet, I know
that geometric algebra is a powerful tool for understanding, as it unifies many
different concepts in linear algebra and geometry. For example, the bivector is
a nice way to represent a plane, which, without geometric algebra, is a distinct
concept that I have to represent using a normal and a point (and extent).

There are some things I really like about geometric algebra, as it relates to
writing a physically based renderer:

1. nicely represents 2D surface interaction, a bivector is a "surfel".
2. nicely represents 3D volume interaction, a trivector is a "voxel".

I've read a lot of different sources and papers, such as John Vince's book on
 Geometric Algebra, and this souce:
https://slehar.wordpress.com/2014/03/18/clifford-algebra-a-visual-introduction/comment-page-1/

and numerous other textbooks and papers.

My understanding is that the way we deal with the frustration of the cross
product and dot product, and the difference between points and vectors, is to
use representations like homogeneous coordinates, which unify both operations
under a single operation of a matrix multiplication.

The Gibbs notation, used by Heaviside to fomulate Maxwell's equations, was
later extended by Hamilton who introduced the quaternions, which among other
things, allowed for composition of rotations.

I am interested in learning geometric algebra primarily because:

1. I want to understand Maxwell's equations in its one equation form.
   I have always struggled to intuitively grasp the four equations.
2. I want to understand intuitively using points, lines, planes, and volumes,
   as primitives, rather than vectors and matrices. By working through each
   step to understand the alternative representation, I hope to gain a deeper
   understanding of the underlying concepts. I want to later reflect and compar
   to better understand the differences.
3. I want to explore in a sandbox way to see if there are practical applications
   and whether I see insights that I wouldn't have seen otherwise.

I also think when I read some of the pros and cons of both approaches online,
that the benefits and drawbacks are sometimes exaggerated. I think that
ultimately, the geometric algebra approach may be a more intuitive way to think
about geometry, because of its powerful set of primitives that extend to
volume rendering applications.

The primary con I see, if I planned to open source this renderer, is that the
use of geometric algebra is simply not going to be familiar to anyone else.

Ultimately, the reason for thinking about it is that I want to understand how
the model translates, so that I can generalize my conceptual model of the
renderer to be more flexible and powerful. I want to be able to recognize
patterns in Gibbs notation and Hamilton's quaternions, and see how they relate
to the concepts I've learned about in geometric algebra.

I hope that the additional insight can help me with future projects such as
rendering point clouds, volume rendering, or other image based rendering task,
or even reverse engineering point clouds to get CAD models. That might be
because I've had so much practice with the conceptual model that I simply get
better over time.

I anticipate it could also help me with visualizing VR applications, or even
neural network computations, as they tend to have geometic algebra analogs.

# 2024-11-14 commands for tracing and inspecting the trace in gotraceui.

I read a lot about the trace package and how it can be used to trace tasks,
regions, and supports event logging with context.

I used this command to download and run gotraceui, a popular GUI tool for
inspecting Go traces:

```powershell
go run -ldflags="-H windowsgui" honnef.co/go/gotraceui/cmd/gotraceui@latest
```

I also implemented my very first fuzz test in Go, and I think it is kind of
neat the idea that you could use it to help generate vectors and shapes
for the raytracer. In a way, testing a renderer is a fantastic use case for
fuzzing because you want the system to work for a wide range of potentially
matheatically degenerate cases.

The default fuzz test time is apparently 1 second and they encourage using
short times for fuzz tests. To run my fuzz test for longer, I use this command:

```powershell
 go test -fuzz=Fuzz -fuzztime=10s
```

Where `-fuzz` must match one fuzz test pattern, and `-fuzztime` is the time
to run the fuzz test for, which defaults to one second.

Today, oh boy, I coded a ton and did a big cleanup and refactor and improved docs.
The biggest things I implemented today are more sophisticated .OBJ parsing
and I also got started on eventually merging my vector types into an f64
package so that I can share it with other projects that won't need to import
phys just to access basic vector types.

I've been very excited by the prospect of getting more advanced OBJ loading
working so that I can render the surfaces and color textures I generated
using my Creality Scan Otter 3D scanner. I really like that scanner by the way.


# 2024-11-14 .obj object file is the objective

The process of writing a .obj parser forced me to confront some details in the
math model I had avoided in my earlier version of raytracee.

I added detailed execution tracing which is amazing for debuggin what happens
in every ray interaction event. I did find it slower however, so I thought that
maybe I would need to add a sampling flag or something.

Instead, while reading more source code on the go/x/ packages, I found that
there is already a really nifty way to do rate limiting on log messages
with the `Sometimes` type provided in go/x/time/rate. I love the Sometimes.Do
method, it is so simple and elegant. I'm also a big fan of how the rate limit
parameters are mixable, to cover all of my use cases.

```go
package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	s := rate.Sometimes{
		First:    2,
		Every:    2,
		Interval: 2 * time.Second,
	}
	s.Do(func() { fmt.Println("1 (First:2)") })
	s.Do(func() { fmt.Println("2 (First:2)") })
	s.Do(func() { fmt.Println("3 (Every:2)") })
	time.Sleep(2 * time.Second)
	s.Do(func() { fmt.Println("4 (Interval)") })
	s.Do(func() { fmt.Println("5 (Every:2)") })
	s.Do(func() { fmt.Println("6") })
}
```

The rate package was added in 2022 according to the package source copyright header.

I also found some excellent free 3D models from NASA hosted on GitHub here:
https://github.com/nasa/NASA-3D-Resources/

I thought some of the heightmaps and asteroid textures might be useful, and it
also includes really awesome models cube sat, curiosity rover, and other
cool spacecraft and satellites.

My mission right now is to successfully load and render both a surface and
color 3D texture from a .obj+mtl file. I think I'm getting close to that goal.
I am now able to load a cleaned up .obj file, but I am still working out how
to map the texture coordinates to the surface. I hope to overcome these last
significant hurdles and get my awesome 3D scans into my renderer!

I also need to refactor the obj loader so that it can handle obj files with
multiple objects in them. This means that an obj file roughly corresponds tp
a `phys.Scene` in my renderer.

I've been also really keen to try plotting a Sankey sort of diagram to show
the radiance energy budget for a ray, starting with the primary ray, just like
the energy budget diagrams in atmospheric science showing the flow of energy
incident from the science as it is absorbed, scattered, reflected, and emitted.

I also found an excellent resource summarizing different file formats here:
https://people.sc.fsu.edu/~jburkardt/data/data.html

This guy really inspires me in terms of data organization and curation and
distillation of knowledge. I think I could learn a lot from his approach to
organizing his information. He has so many categorized 3d file examples that
are so helpful to me when working on my renderer.


# 2024-11-15 continuing the struggle of .obj texture support.

My mission is to render a .obj and .mtl (and .png texture map) for one of the
3D scan models that I collected with my Creality Scan Otter 3D scanner.

So far I've implemented a .obj filesystem loader that can load an object file
and any of the materials and textures associated with it. Awesome.

I've also implemented a converter between the object file representation and
my scene entity representation. This part works for loading surfaces from .obj
files with only a single mesh inside.

I read that .obj can support multiple objects and materials in a single file,
and if that is the case, then it really maps more to a Scene than to a single
phys.Entity.

On a side topic, I'm really liking how the phys package is turning out.
The process of building a raytracer is truly such a rewarding and challenging
experience that I know will form the basis for most of my future imaging and
graphics related projects.

There are a few things to do for my obj package:

1. Add testdata for a range of obj/mtl input scenarios, multiple objects, etc.
2. Support groups which naturally map to phys.Group. Similar to Three.js loader.
3. Improve error messages with more surrounding context like ParseError.

Test scenarios:

1. single object, single material, single texture
2. multiple objects, single material that is inherited by the others
3. multiple objects, multiple materials, multiple textures

Other possible sanity checks:

1. verify indices reference existing vertices, texture coords, normals.
2. verify number of vertices, texture coords, normals match up.
3. material properly parsed and associated with faces.
4. performance time budget

One of the best development references for OBJ loading is the Three.js loader
for OBJ files: OBJLoader.js
https://github.com/mrdoob/three.js/blob/master/examples/jsm/loaders/OBJLoader.js

Although a very simple format to get started with on simple models, there is a
surprising amount to learn about from the OBJ format, such as the different
ways to represent normals, texture coordinates, and materials. I found the
process of learning to parse the file from scratch to be very rewarding.
I think to an extend this is necessary, since the understanding gained from
developing the loader is essential to understanding how it interfaces with the
phys package and the related graphics types.

With a variety of small performance improvements, I was able to refactor the obj
package from 77 MB/s to 94 MB/s, over a 20% improvement! I found that it took a
long time to load some of the 3D scanned .obj files with over a million vertices
and faces, so I wanted to make sure that the loader was as fast as possible.

I could go crazy and make a parallel loader that optimistically loads large
amounts of vertices but then reconstructs the scene in a single goroutine
afterwards, correcting and possibly editing the scene to correct errors caused
by optimistic concurrency. I think that would be a fun project to work on, but
I don't plan to do that right now. Still, it does take several annoying seconds
to load a detailed OBJ file, and if that continues to cause headache I could
also investigate other options such as loading it once and then caching it with
gob or another serialization format.

I'm still thinking about ways to improve the package API for obj files.
I want something idiomatic and uses function names simnilar to the standard
library packages.

The best resource on the obj file is here:

https://paulbourke.net/dataformats/obj/

# 2024-11-15 fuzz testing a physically based renderer

I found fuzz testing to be a suprisingly interesting application for my `phys`
physically based renderer package.

Use cases:

1. Generate random triangles and shapes for collisions.
2. Generate random rays for testing intersections.
3. OBJ file: generate random vertices and faces for loading.

My favourite is thinking about how to design the rendering and surface code
to work with fuzz testing, because it is able to generate so many more possible
cases for rays that my limited test scenes cannot. It can discover, edge cases,
literally edge cases, where rays are just barely touching the surface, or
where the surface is just barely visible to the camera. All of these things
build confidence in the correctness of the renderer. At the same time, I'm
impressed by how clean the Fuzz testing implementation in Go is. It really
feels like a natural part of the language and tooling ecosystem.

The example below shows two fuzz tests I implemented.
The first fuzz test generates random elements for an obj file and verifies
that the obj parser handles errors gracefully, without panic or crash.

The second fuzz test it really cool. It generates randomly orientied rays
and collides them with a triangle. It then verifies that the collision
point and normal are correct. This is a really powerful test because it
can generate so many more cases than I could ever write by hand.

As soon as I wrote the second fuzz test, it immediately found a failing test
case in my triangle collision code. I was so impressed by how quickly it found
the bug, and how easy it was to write the test. I'm curious to see what other
applications I find for it while writing my phys renderer.

The example below shows running the first fuzz test successfully and the second
fuzz test failing, which produced really interesting output.

```
PS C:\Users\scott\github\slam\code\raytrace\obj>  go test -fuzz=Fuzz -fuzztime=10s
fuzz: elapsed: 0s, gathering baseline coverage: 0/130 completed
fuzz: elapsed: 0s, gathering baseline coverage: 130/130 completed, now fuzzing with 6 workers
fuzz: elapsed: 3s, execs: 305218 (101727/sec), new interesting: 93 (total: 223)
fuzz: elapsed: 6s, execs: 651671 (115099/sec), new interesting: 146 (total: 276)
fuzz: elapsed: 9s, execs: 976689 (108619/sec), new interesting: 181 (total: 311)
fuzz: elapsed: 11s, execs: 1045404 (33462/sec), new interesting: 183 (total: 313)
PASS
ok      github.com/scottlawsonbc/slam/code/photon/raytrace/obj 11.245s
PS C:\Users\scott\github\slam\code\raytrace\obj> cd ../phys
PS C:\Users\scott\github\slam\code\raytrace\phys>  go test -fuzz=Fuzz -fuzztime=10s
--- FAIL: TestTriangleCollideEdgeCases (0.00s)
    --- FAIL: TestTriangleCollideEdgeCases/Ray_grazes_the_triangle_(t_==_tmin) (0.00s)
        shape_triangle_test.go:416: Expected hit: false, got: true
        shape_triangle_test.go:420: Expected intersection point: (0, 0, 0), got: (0.5, 0.5, 0)
        shape_triangle_test.go:423: Expected normal: (0, 0, 0), got: (0, 0, 1)
    --- FAIL: TestTriangleCollideEdgeCases/Ray_starts_on_the_triangle_and_points_away (0.00s)
        shape_triangle_test.go:416: Expected hit: false, got: true
        shape_triangle_test.go:420: Expected intersection point: (0, 0, 0), got: (0.25, 0.25, 0)
        shape_triangle_test.go:423: Expected normal: (0, 0, 0), got: (0, 0, 1)
FAIL
exit status 1
FAIL    github.com/scottlawsonbc/slam/code/photon/raytrace/phys        0.162s
PS C:\Users\scott\github\slam\code\raytrace\phys>
```

I haven't figured out what the source of the bug is yet, will investigate more.

# 2024-11-16 sudden moment of insight about the fundamentalness of triangles

I know about triangles. I know how to implement them in code. I've been using
them for years in my raytracer code. I've also known about shaders, and I've
read so many things about shaders and what they do, but I've never really been
satisfied with my understanding of them. There was always something missing,
an energy gap between seeing them in games and applications and seeing the
mathematical representation of them in code.

My mission: load and render a 3D scan that I captured with my 3D scanner.
My scanner outputs .obj and .mtl files that reference 4096x4096 .png texture
maps.

I was relatively quickly able to load the .obj file as a shape without texture,
as that's something which my renderer has supported for a long time.

I ran into a roadblock when it came to loading the texture map.
After a lot of reading into the whole process of texture mapping, I suddenly
realized the purpose of having a mesh data structure distinct from that of a
group of triangles. I suddenly realized in all of these different ways why
that extra mapping of to UV in the triangle is so fundamental to how rendering
works, it is critical that the triangle be able to communicate to the material
which part of the texture map it should use. It really is fundamental because
of the separation between material and shape. Instead of tightly coupling the
material and shape, we can instead add the extra UV mapping to the triangle.

In summary, I was never satisfied with my understanding of other people's
triangle code because it always seemed to have an extra set of data that I
didn't understand the purpose of, since my renderer didn't support textures.
Now that I'm trying to add texture support, I suddenly see the purpose of that
extra data, and the moment of insight is so satisfying.


# 2024-11-16 thoughts while reading through some of the gonum library

I had heard about gonum for a couple of years now, at least, but I had never
read through the source code until today. I'm really inspired by the quality
of the code and how I can just click into the directory and read and understand
the code.

A few files that stuck out to me as really interesting:

fd - diff.go - finite difference derivatives
comment: I like the Formula type and how it can represent different methods
for calculating derivatives. I also like the way the code is organized into
packages for more relatively narrow aspects of mathematics.

stat directory
comment: I love this collection of tested statistics functions.
As I get into more sampling and probability distribution methods in my phys
renderer, I think I will want to come back here and use some of the methods.

triangle.go - I love that there is a tested method here for a triangle and the
core mathematical operations that are needed to work with it.

gonum/blas - I can see immediately why something like this has potential as a basic
foundation for a lot of the more complicated linear algebra methods that I may
encounter in future parts of the project, such as continuous surfaces or volume
rendering.

gonum/unit - Really interesting to see how the different units are defined and how
they are used in the code. I like how each unit is in a different file.
I thought it was also nice that they used a go idiomatic consistent interface
naming convention, for example, the unit Current implements "Currenter".

r2 and r3 - I like how this is organized into separate packages for 2D and 3D.
I don't think I would have thought about doing that.
I am so glad that I can use this as a reference for how to organize my own code,
especially the geometric types which are newest to me.
Another thing I like about this implementation is that it is a different take
on how I've been naming my vectors, such as `Vec2` and `Vec3`. I like the
idea of that information being in the package name instead of the type name.
I like that this frees us to use more context specific names for the types.

gonum/r3/box.go - I like the shorter name and optimized methods for the box type.
I think I may adopt something like this.

gonum/r3/mat_unsafe.go - Wait, what is this? A matrix implementation with unsafe
pointers! That's really interesting because I don't think they would have made
that for no reason, and it makes me wonder whether I would get a boost from
using the unsafe pointers in my own rendering code. More of a curiosity than
anything else. I suspect the matrix operations can feel the function call tax
and that's why they use unsafe pointers.

gonum/mat/io.go - I like how they have methods for marshalling and unmarshalling
the matrix data. This is a good idea.

gonum/interp/interp.go - I like the Fitter and Predictor interfaces. Really neat
way to model the different interpolation methods.

gonum/optimize package - AMAZING that I can see these different optimization
methods, that's so freaking cool. They have a localMethod, lineSearcher,
NextDirectioner, StepSizer, Recorder, Statuser. What a nice breakdown of the
different methods. This package has gradient descent and a wide range of other
methods that I've used in numpy and matlab. Such cool functions here.
I'm sure there are multiple ways in which I could use this to do certain render
related things like minimum variance sampling, or optimizing texture and models
or bounding volume hierarchies.

gonum/internal/asm - Cool collection of optimized assembly routines. I am curious
what sort of performance boost I could get from using these in my own code.

gonum/dsp - Really cool. It provides a fourier transform implementation and
some windowing functions. If I were to make audioled using Go, this is the
sort of package I would look at first.

gonum/mathext - For a bunch of nice math functions that aren't part of the
standard gonum package, but which show up a lot in fields like physics and
statistics.

gonum/mat package - Everything you could possibly want in a matrix library.

gonum/graph package - I like having these graph algorithms available to me, I should
read more on this later. What a fascinating codebase to read through.
I like the Graph interface and I wonder if that also applies to my ray tree in
the renderer.

gonum/graph/traverse/traverse.go - I like how the traversal works on such a
generic interface, like with BreadthFirst and DepthFirst.

I hesitate to use the gonum code directly for the time being, as I want the
phys package to only have code that I fully understand. However, in the future,
I think I will want to reference or adapt some of it.

Overall, I am really impressed by the clean, understandable, and well organized
code in the gonum library. I will either use this package in the future or learn
a lot from it even without using it.

In many ways it is helpful to jog my own ideas by looking at their code and
seeing how they implement things like tests, or represent varying methods
for the same underlying fundamental operations, or how they organize their
code into packages and files.


Unrelated to gonum, but similar to the fogelman pt pathtracer, there is another
pbr package which uses some alternative representations.

https://github.com/hunterloftis/pbr/blob/master/pkg/bsdf/microfacet.go

This one uses different packages to organize the various separable aspects.

https://github.com/hunterloftis/pbr/blob/master/pkg/camera/slr.go

pbr/pkg/camera/slr.go - I like how this implements an SLR camera, which should
be a good reference for my own camera implementation.


pbr/pkg/farm/server.go - Looks like it implements a basic HTTP render interface
which is pretty cool.



# 2024-11-16 why didn't I change my vscode path separator on windows before?

I have been bothered by the fact that, on Windows, vscode uses the backslash
as the path separator, while Windows conventially uses the backslash.

Windows supports the backlash as well, but it is not the default.

Today it finally occured to me that maybe there is a setting in vscode that I
can change so that it always uses the forward slash as the path separator.
And, well, there is!

The setting is called:

`explorer.copyRelativePathSeparator` and I set it to "/" instead of the default
of "auto".

Why the heck didn't I do this before? I guess I just accepted it and never
thought about it, or it always came up while I was in the middle of something
else.


# 2024-11-16 triangle, why you no work? How can triangles be so difficult

My test capsule object does not seem to be working. The craziest part is that
I don't get an obvious texture mapping mistake, it seems almost random and
glitchy, and that part is really crazy to understand.

Example of the problem:

code/raytrace/scene/scan/out/out_20241115_185926.png

I'm not sure what the problem is, and at the moment, the problem isn't very
easy to debug, so I'm trying to brainstorm both how to move forward on the
texture mapping problem, what things I can do to make it easier to debug this
and other things, and how I can simplify the test case to better understand
where the problem is.

I think afterwards I may go back and make significant more changes to the
triangle code, but at the moment, I'd like to get textures working at all.

First, I will identify a simple test case to do sanity checks on reported
texture coordinates and normals.

I will also implement a UV shader to visualize the texture coordinates on the
surface of the object. I may be able to spot sudden discontinuities or other
problems that way.

# 2024-11-16 texture resources

Some good texture 3D model assets here (mesh and texture):

https://texturedmesh.isti.cnr.it/browse

Lots of good statue and other still 3D models.

Once I get textures working, I should come back here and try to render some
of these, since they are high quality and have a lot of detail, would make
for some impressive demo scenes.


# 2024-11-17 more go tooling exploration

golang/mobile/app/internal/apptest - Really cool Comm struct and how it is
passed a bufio.Scanner and so it doesn't need to worry about how to unpack
or pack program specific data.
https://github.com/golang/mobile/blob/master/app/internal/apptest/apptest.go

golang/mobile/gl - This is really cool, I should definitely use this as a
reference if I decide to make my phys renderer support opengl background.
https://github.com/golang/mobile/tree/master/gl The example they have seems
to work on Windows and Linux so that's awesome.

I really like the App interface, and how clean the PublishResult is, and how
all communication comes in the form of decoupled events.
```go
type App interface {
	// Events returns the events channel. It carries events from the system to
	// the app. The type of such events include:
	//  - lifecycle.Event
	//  - mouse.Event
	//  - paint.Event
	//  - size.Event
	//  - touch.Event
	// from the golang.org/x/mobile/event/etc packages. Other packages may
	// define other event types that are carried on this channel.
	Events() <-chan interface{}

	// Send sends an event on the events channel. It does not block.
	Send(event interface{})

	// Publish flushes any pending drawing commands, such as OpenGL calls, and
	// swaps the back buffer to the screen.
	Publish() PublishResult

	// Filter calls each registered event filter function in sequence.
	Filter(event interface{}) interface{}

	// RegisterFilter registers a event filter function to be called by Filter. The
	// function can return a different event, or return nil to consume the event,
	// but the function can also return its argument unchanged, where its purpose
	// is to trigger a side effect rather than modify the event.
	RegisterFilter(f func(interface{}) interface{})
}
```


# 2024-11-18 what a unity material file looks like

I went through my old game asset collection and in the `PolygonSciFiSpace` pack
I found an example `.mat` file:

```yaml
%YAML 1.1
%TAG !u! tag:unity3d.com,2011:
--- !u!21 &2100000
Material:
  serializedVersion: 6
  m_ObjectHideFlags: 0
  m_PrefabParentObject: {fileID: 0}
  m_PrefabInternal: {fileID: 0}
  m_Name: PolygonScifiSpace_Material_04_A
  m_Shader: {fileID: 46, guid: 0000000000000000f000000000000000, type: 0}
  m_ShaderKeywords: _EMISSION
  m_LightmapFlags: 1
  m_EnableInstancingVariants: 0
  m_CustomRenderQueue: -1
  stringTagMap: {}
  disabledShaderPasses: []
  m_SavedProperties:
    serializedVersion: 3
    m_TexEnvs:
    - _BumpMap:
        m_Texture: {fileID: 0}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _DetailAlbedoMap:
        m_Texture: {fileID: 0}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _DetailMask:
        m_Texture: {fileID: 0}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _DetailNormalMap:
        m_Texture: {fileID: 0}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _EmissionMap:
        m_Texture: {fileID: 2800000, guid: e29a7a5f451453f4f978120a03c0f52b, type: 3}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _MainTex:
        m_Texture: {fileID: 2800000, guid: 387a11e4dfd740f4095be2986edfa42d, type: 3}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _MetallicGlossMap:
        m_Texture: {fileID: 0}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _OcclusionMap:
        m_Texture: {fileID: 0}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    - _ParallaxMap:
        m_Texture: {fileID: 0}
        m_Scale: {x: 1, y: 1}
        m_Offset: {x: 0, y: 0}
    m_Floats:
    - _BumpScale: 1
    - _Cutoff: 0.5
    - _DetailNormalMapScale: 1
    - _DstBlend: 0
    - _GlossMapScale: 1
    - _Glossiness: 0
    - _GlossyReflections: 1
    - _Metallic: 0
    - _Mode: 0
    - _OcclusionStrength: 1
    - _Parallax: 0.02
    - _SmoothnessTextureChannel: 0
    - _SpecularHighlights: 1
    - _SrcBlend: 1
    - _UVSec: 0
    - _ZWrite: 1
    m_Colors:
    - _Color: {r: 1, g: 1, b: 1, a: 1}
    - _EmissionColor: {r: 0.5, g: 0.5, b: 0.5, a: 1}
```

Really interesting to compare that to my own `phys` material representation,
and to the `.mtl` file used in the `.obj` file format. The unity one has far
more fields, and I may want to reference this in the future when it comes to
adding more fields to my material. I'm sure everything in the unity material
has a purpose, so it's worth at least thinking about.


# 2024-11-18 scott's texture nightmare continues

I'm trying to figure out why I'm unable to render textures properly for a 3D
model that I scanned. Meshlab, however, can load the file perfectly.

Even more mysterious, showing the texture mapping in meshlab and comparing it
to what I see in the texture image file, it's slightly different. There are
some faces shown in wireframe in meshlab that are shown in solid color in the
texture .png file.

To investigate what's going on, I 3D scanned my isopropyl spray bottle, as it
has a nice and smooth surface with high texture contrast. When I rendered the
image,


However, this simple pirate bottle renders without an issue:

<!-- slam\code\raytrace\scene\scan\out\out_20241117_162243.png -->

<img src="./scene\scan\out\out_20241117_162243.png" width="200">


The strangest part is how the texture shows up as partially correct but with
obvious errors. I tried inverting my texture coordinates, which made changes
to the texture mapping, but the mapping is still incorrect.

I'm also learning more about vertices and how they are associated with texture
attributes, but more generally, any attribute that can be associated with a
vertex.

I'm not sure if there are multiple conventions for texture coordinates, but I
can't figure out why my 3D scan is not texture mapping properly. Is it possible
that I'm using a different orientation for my texture coordinates than the
texture image file?

I'm using a per-vertex texture coordinate.

I found a bunch of 3D assets in my polygon sci-fi space pack asset!
The files are .fbx with an associated .fbx.meta file.

Example

```yaml
fileFormatVersion: 2
guid: cf8777e2dce11c0458b21b134bbc772e
timeCreated: 1549333825
licenseType: Store
ModelImporter:
  serializedVersion: 19
  fileIDToRecycleName:
    100000: //RootNode
    400000: //RootNode
    2300000: //RootNode
    3300000: //RootNode
    4300000: SM_Tunnel_Mesh
    9500000: //RootNode
  materials:
    importMaterials: 0
    materialName: 0
    materialSearch: 1
  animations:
    legacyGenerateAnimations: 4
    bakeSimulation: 0
    resampleCurves: 1
    optimizeGameObjects: 0
    motionNodeName:
    rigImportErrors:
    rigImportWarnings:
    animationImportErrors:
    animationImportWarnings:
    animationRetargetingWarnings:
    animationDoRetargetingWarnings: 0
    animationCompression: 1
    animationRotationError: 0.5
    animationPositionError: 0.5
    animationScaleError: 0.5
    animationWrapMode: 0
    extraExposedTransformPaths: []
    clipAnimations: []
    isReadable: 1
  meshes:
    lODScreenPercentages: []
    globalScale: 1
    meshCompression: 0
    addColliders: 0
    importBlendShapes: 1
    swapUVChannels: 0
    generateSecondaryUV: 0
    useFileUnits: 1
    optimizeMeshForGPU: 1
    keepQuads: 0
    weldVertices: 1
    secondaryUVAngleDistortion: 8
    secondaryUVAreaDistortion: 15.000001
    secondaryUVHardAngle: 88
    secondaryUVPackMargin: 4
    useFileScale: 1
  tangentSpace:
    normalSmoothAngle: 60
    normalImportMode: 0
    tangentImportMode: 3
  importAnimation: 0
  copyAvatar: 0
  humanDescription:
    serializedVersion: 2
    human: []
    skeleton: []
    armTwist: 0.5
    foreArmTwist: 0.5
    upperLegTwist: 0.5
    legTwist: 0.5
    armStretch: 0.05
    legStretch: 0.05
    feetSpacing: 0
    rootMotionBoneName:
    rootMotionBoneRotation: {x: 0, y: 0, z: 0, w: 1}
    hasTranslationDoF: 0
    hasExtraRoot: 0
    skeletonHasParents: 1
  lastHumanDescriptionAvatarSource: {instanceID: 0}
  animationType: 0
  humanoidOversampling: 1
  additionalBone: 0
  userData:
  assetBundleName:
  assetBundleVariant:
```

I realized just now that I actually have not tested my mesh shape, I have only
tested my triangle and quad shapes. I should add a sphere to my texture scene
as well.

I think I realized that since I never tested the mesh model, and I never used
my new TriangleUV type, that I actually haven't really tested my renderer.

# 2024-11-18 more cool finds in go packages
Looking through go packages and saved a few interesting finds:

gonum/spatial/r3 - I never would have thought to write a handler like this, but it is very cool!
I like how it abstracts the field as a func(Vec) Vec. I wonder if I could use this
to render some really interesting mathematical or physical surfaces in my phys renderer.
```go
// Divergence returns the divergence of the vector field at the point p,
// approximated using finite differences with the given step sizes.
func Divergence(p, step Vec, field func(Vec) Vec) float64 {
	sx := Vec{X: step.X}
	divx := (field(Add(p, sx)).X - field(Sub(p, sx)).X) / step.X
	sy := Vec{Y: step.Y}
	divy := (field(Add(p, sy)).Y - field(Sub(p, sy)).Y) / step.Y
	sz := Vec{Z: step.Z}
	divz := (field(Add(p, sz)).Z - field(Sub(p, sz)).Z) / step.Z
	return 0.5 * (divx + divy + divz)
}

// Gradient returns the gradient of the scalar field at the point p,
// approximated using finite differences with the given step sizes.
func Gradient(p, step Vec, field func(Vec) float64) Vec {
	dx := Vec{X: step.X}
	dy := Vec{Y: step.Y}
	dz := Vec{Z: step.Z}
	return Vec{
		X: (field(Add(p, dx)) - field(Sub(p, dx))) / (2 * step.X),
		Y: (field(Add(p, dy)) - field(Sub(p, dy))) / (2 * step.Y),
		Z: (field(Add(p, dz)) - field(Sub(p, dz))) / (2 * step.Z),
	}
}
```

I'm finally moving my math packages to r3 and r2.

gonum/floats/floats.go - I like the methods for adding up []float64, as that
can be a pretty common operation. I don't think I have a use for it at the
moment but I will keep the idea in mind.


# 2024-11-18 texture success!

YAY IT FINALLY WORKS!

I reimplemented my mesh class and finally got texture mapping to work.
This is super exciting!

<img src="./scene/scan/out/out_20241119_015515.gif" width="512">


# 2024-11-20 basic texture editing with 3d tools

I've been using meshlab and now just downloaded blender. I have a few use
cases that I need a tool to help me with:

1. I want to be able to edit the texture of a 3D scan mode
2. I want to be able to delete vertices from a 3D scan model:
  1. to fix holes.
  2. to fix 3D scan errors
3. I want to be able to smooth a 3D scan model surface
4. I want to be able to paint or retouch a 3D scan model

What are the fundamental mesh operations I can do with the tools?
I like how blender lays out everything in different tabs.

The "Modeling" tab has basic mesh operations like extrude, inset, bevel, etc.

I want to write down a list of my concrete mesh editing use cases first, then
I will be able to look up how to do them in my mesh tools like blender or
meshlab.

1. Scott select face and paint texture over 3D scan model.
2. Scott select vertices and align mesh to plane.
3. Scott select

# 2024-11-19 molecular dynamics simulation and protein aggregation

I learned about GROMACS, a molecular dynamics simulation package that is used
to simulate protein folding and aggregation. I'm curious to learn more about
how these things are simulated.
http://www.mdtutorials.com/gmx/


# 2024-11-19 starting on glTF journey

I'm on a quest to implement a gltf package so that I can load and export to
.glTF files. I really like this as a 3D format and I want to start implementing
it earlier so that `phys` is neatly modeled after gltf and is easy to convert.

I'm documenting some of notes and misc thoughts about gltf as I go:

Interesting that Scene has a list of nodes, and that a gltf asset may contain
multiple scenes.

Also that when an asset has no scenes, then:

> A glTF asset that does not contain any scenes SHOULD be treated as a library of individual entities such as materials or meshes.

https://registry.khronos.org/glTF/specs/2.0/glTF-2.0.html#concepts

The nodes that scenes contain must be the root level nodes of the scene.
Nodes may be a hierarchy, but only root level notes are referenced by scenes.



# 2024-11-19 tracing the tracing

In code/raytrace/scene/scan I ran go pprof and go this output for `top10`:

```
PS C:\Users\scott\github\slam\code\raytrace\scene\scan> go run . -cpuprofile "cpuprof.out"
Writing CPU profile to cpuprof.out
Walking &{../../../../3d/scan/ linear-stage-controller-marker-simplified.obj} msg=main.modelFS
2024/11/19 20:43:38 | linear-stage-controller-marker-simplified.mtl 268 230d0a23205761766566726f6e74206d d47b867a6bcaf9d27f57ca58eaee6b55
2024/11/19 20:43:38 | linear-stage-controller-marker-simplified.obj 21647257 232323230d0a230d0a23204f424a2046 97004f1b2f6cf29e16707367bdeb9ebf
2024/11/19 20:43:38 | linear-stage-controller-marker-simplified.png 9075558 89504e470d0a1a0a0000000d49484452 6d9191d4cc9773471d50ac4318539a4f
2024/11/19 20:43:38 loading texture linear-stage-controller-marker-simplified.png
got 1 nodes
node 0 bounds {(-100.1248169, -100.1051483, -64.7324371) (100.1248169, 100.1051483, 64.7324371)}
Rendering scene 0 with camera {LookFrom:(1e+11, 0, 1e+11) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000 nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 20:43:46 Rendered 512x512
        RenderTime: 3.1655931s (12.075µs per pixel)
        IntersectionTests: 22534400
        TotalRays: 22534400
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 11545618 (51.2%)
Rendering scene 1 with camera {LookFrom:(6.1232339957367576e-06, 1e+11, 2e+10) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000 nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 20:43:48 Rendered 512x512
        RenderTime: 5.8705924s (22.394µs per pixel)
        IntersectionTests: 45054464
        TotalRays: 45054464
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 27000753 (59.9%)
Rendering scene 2 with camera {LookFrom:(-1e+11, 1.2246467991473515e-05, 9.999999999999998e+10) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000
nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 20:43:52 Rendered 512x512
        RenderTime: 9.1230272s (34.801µs per pixel)
        IntersectionTests: 67627520
        TotalRays: 67627520
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 37696114 (55.7%)
Rendering scene 3 with camera {LookFrom:(-1.8369701987210274e-05, -1e+11, 1.8e+11) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000 nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 20:43:55 Rendered 512x512
        RenderTime: 12.2206475s (46.618µs per pixel)
        IntersectionTests: 90141952
        TotalRays: 90141952
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 48127037 (53.4%)
2024/11/19 20:43:55 stored 115032583 bytes of trace data 109 MB
Saved to ./out/out_20241119_204355.gif
2024/11/19 20:43:55 Saved to ./out/out_20241119_204355.gif (109 MB)
```

```
(pprof) top10
Showing nodes accounting for 37140ms, 74.28% of 50000ms total
Dropped 321 nodes (cum <= 250ms)
Showing top 10 nodes out of 95
      flat  flat%   sum%        cum   cum%
   14570ms 29.14% 29.14%    28920ms 57.84%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.Group.Collide
    7510ms 15.02% 44.16%    10660ms 21.32%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.Face.Collide
    4660ms  9.32% 53.48%    44070ms 88.14%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.Render.func2
    3540ms  7.08% 60.56%     3540ms  7.08%  runtime.duffcopy
    1900ms  3.80% 64.36%    33830ms 67.66%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.(*BVH).Collide
    1500ms  3.00% 67.36%     1520ms  3.04%  github.com/scottlawsonbc/slam/code/photon/raytrace/r3.Vec.Dot (inline)
    1280ms  2.56% 69.92%     2970ms  5.94%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.AABB.hit
     820ms  1.64% 71.56%      830ms  1.66%  github.com/scottlawsonbc/slam/code/photon/raytrace/r3.Point.Sub (inline)
     700ms  1.40% 72.96%     1590ms  3.18%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.NewBVH.func1
     660ms  1.32% 74.28%      660ms  1.32%  runtime.asyncPreempt
(pprof)
```

```
i 100 timePerSplit 2134157us progress 0.050500% opsPerSec 46.856909
i 200 timePerSplit 2243641us progress 0.100501% opsPerSec 44.570410
i 300 timePerSplit 2194136us progress 0.150501% opsPerSec 45.576026
i 400 timePerSplit 2124130us progress 0.200501% opsPerSec 47.078098
```

This is crazy, when I switch to surface area heuristic, the render time is
crazy slow! The initialization time is insane.

With this awful speed I now implemented a much faster binned SAH implementation.
When I tried SAH without binning, it was still very slow.

It is AWESOME! The new binned SAH BVH shape is 5us per pixel down from 12, a
more than 2x speedup. I'm so happy with this result. I'm actually shocked
at how fast it is, especially since compared to the previous implementation
which was already a bounding volume hierarchy.

```
Rendering scene 0 with camera {LookFrom:(1e+11, 0, 1e+11) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000 nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 23:02:50 Rendered 512x512
        RenderTime: 1.3174364s (5.025µs per pixel)
        IntersectionTests: 22509824
        TotalRays: 22509824
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 11514513 (51.2%)
Rendering scene 1 with camera {LookFrom:(6.1232339957367576e-06, 1e+11, 2e+10) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000 nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 23:02:51 Rendered 512x512
        RenderTime: 2.40844s (9.187µs per pixel)
        IntersectionTests: 45051904
        TotalRays: 45051904
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 26997031 (59.9%)
Rendering scene 2 with camera {LookFrom:(-1e+11, 1.2246467991473515e-05, 9.999999999999998e+10) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000
nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 23:02:52 Rendered 512x512
        RenderTime: 3.4195861s (13.044µs per pixel)
        IntersectionTests: 67582464
        TotalRays: 67582464
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 37638802 (55.7%)
Rendering scene 3 with camera {LookFrom:(-1.8369701987210274e-05, -1e+11, 1.8e+11) LookAt:(0, 0, 0) VUp:(0, 0, 1) FOVHeight:250.000000 nm FOVWidth:250.000000 nm}
Rendering: 100% completete
2024/11/19 23:02:53 Rendered 512x512
        RenderTime: 4.5191553s (17.239µs per pixel)
        IntersectionTests: 90112512
        TotalRays: 90112512
        RaysExceedingDepth: 0 (0.0%)
        RaysLeftScene: 48067646 (53.3%)
2024/11/19 23:02:53 stored 114329495 bytes of trace data 109 MB
Saved to ./out/out_20241119_230253.gif
2024/11/19 23:02:53 Saved to ./out/out_20241119_230253.gif (109 MB)
2024/11/19 23:02:53 Saved to ./scan.png
2024/11/19 23:02:54 Saved to ./trace.out
```

```
(pprof) top10
Showing nodes accounting for 11820ms, 56.91% of 20770ms total
Dropped 290 nodes (cum <= 103.85ms)
Showing top 10 nodes out of 131
      flat  flat%   sum%        cum   cum%
    4710ms 22.68% 22.68%    16610ms 79.97%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.Render.func2
    1890ms  9.10% 31.78%     6670ms 32.11%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.(*BVH).Collide
    1320ms  6.36% 38.13%     3290ms 15.84%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.AABB.hit
     770ms  3.71% 41.84%     1390ms  6.69%  github.com/scottlawsonbc/slam/code/photon/raytrace/phys.(*Group).Collide
     740ms  3.56% 45.40%      770ms  3.71%  github.com/scottlawsonbc/slam/code/photon/raytrace/r3.Point.Get
     550ms  2.65% 48.05%      550ms  2.65%  runtime._ExternalCode
     520ms  2.50% 50.55%      520ms  2.50%  runtime.stdcall2
     460ms  2.21% 52.77%      460ms  2.21%  runtime.procyield
     450ms  2.17% 54.94%      450ms  2.17%  runtime.memmove
     410ms  1.97% 56.91%      410ms  1.97%  math.archMin
(pprof)
```

This shows that Group.Collide is now a tiny fraction of the time it was before.

For the web interface, I ran this command:

```
go tool pprof -http=localhost:8032 scan cpuprof.out
```


Related: last couple of years has seen updates for profile guided optimization
in Go 1.21. https://go.dev/blog/pgo


# 2024-11-19 more graphics and geometry resources

Came across this site geometric tools which has a collection of math, geometry,
image analysis, and physics source code. It is also documented and has books
about it.

https://www.geometrictools.com/

For example

```

•	AllPairsTriangles:	All-pairs intersection testing between triangles of two meshes (CPU, GPU).
•	IntersectBoxCone:	Test-intersection queries for an aligned/oriented box and a (finite) cone.
•	IntersectBoxCylinder:	Test-intersection queries for an aligned/oriented box and a (finite) cylinder.
•	IntersectBoxSphere:	Test-intersection queries for an aligned/oriented box and a sphere.
•	IntersectConvexPolyhedra:	Find-intersection query for two convex polyhedra.
•	IntersectCylinders:	Test-intersection queries for two finite cylinders.
•	IntersectInfiniteCylinders:	Find-intersection query for two infinite cylinders. The curve of intersection is computed.
•	IntersectLineRectangle:	Find-intersection and test-intersection queries for a line, ray, or segment and a rectangle.
•	IntersectLineTorus:	Find-intersection query for a line and a torus.
•	IntersectPlaneConvexPolyhedron:	Find-intersection query for a plane and a convex polyhedron.
•	IntersectSphereCone:	Test-intersection queries for a sphere and a cone (infinite cone, infinite truncated cone, finite cone, cone frustum).
•	IntersectTriangleBox:	Determine whether a triangle and an oriented box intersect.
•	IntersectTriangleCylinder:	Determine whether a triangle and a cylinder intersect.
•	IntersectTriangles2D:	Test-intersection and find-intersection queries for two triangles in 2D.
•	MovingCircleRectangle:	Illustration of computing first time and point of contact between moving circles and rectangles.
•	MovingSphereBox:	Illustration of computing first time and point of contact between moving spheres and boxes.
•	MovingSphereTriangle:	Illustration of computing first time and point of contact between moving spheres and triangles.
```

https://www.geometrictools.com/Samples/Intersection.html

In particular, I note that `IntersectLineTorus`, `IntersectLineRectangle`,
`IntersectTriangleBox`, `IntersectTriangleCylinder`,


There is a really cool page on scene graphics and I like the BillboardNode,
and MorphControllers (for deformations).

The billboard node always faces the camera, which I think is a really neat
concept that I could use to present some UI elements in 3D space.

https://www.geometrictools.com/Samples/SceneGraphs.html

His github repo is here: https://github.com/davideberly/GeometricTools

I also like the Picking section, about selecting primitives with the cursor.

I want to know how to do surface extraction (SurfaceExtraction section),
from a 3D voxel data set.
I think this is in the neighbourhood of level sets and level curves.

This is the author of the excellent "geometric tools for computer graphics"
book which I downloaded a while ago.


See the paper `extraction of level sets from 3d images` David Eberly.
Level Set Extraction from Gridded 2D and 3D Data
David Eberly, Geometric Tools, R

The problem boils down to doing efficient ray-mesh intersection tests.

Related to the fundamental Moller-Trumbore Algorithm for triangle intersection.

glTF asset concept can hold multiple scenes. If I decide to implement this
too, then I should rename my Scene to Asset and then have a render scene
function for one particular scene.

After spending a bunch of time tracing the tracing, I actually think what
I really need is just pprof.

Types in pprof

```

-flat [default], -cum: Sort entries based on their flat or cumulative value respectively, on text reports.
-functions [default], -filefunctions, -files, -lines, -addresses: Generate the report using the specified granularity.
-noinlines: Attribute inlined functions to their first out-of-line caller. For example, a command like pprof -list foo -noinlines profile.pb.gz can be used to produce the annotated source listing attributing the metrics in the inlined functions to the out-of-line calling line.
-nodecount= int: Maximum number of entries in the report. pprof will only print this many entries and will use heuristics to select which entries to trim.
-focus= regex: Only include samples that include a report entry matching regex.
-ignore= regex: Do not include samples that include a report entry matching regex.
-show_from= regex: Do not show entries above the first one that matches regex.
-show= regex: Only show entries that match regex.
-hide= regex: Do not show entries that match regex.
```

I may want to filter out standard library math functions or any code not
related to the phys package.

The best resource I found for pprof is here:

https://github.com/google/pprof/blob/main/doc/README.md


# 2024-11-20 more glTF notes

There's a good glTF tutorial here;
https://github.com/KhronosGroup/glTF-Tutorials/blob/main/gltfTutorial/gltfTutorial_016_Cameras.md

> Camera instancing and management
> There may be multiple cameras defined in the JSON part of a glTF.
> Each camera may be referred to by multiple nodes.
> Therefore, the cameras as they appear in the glTF asset are really
> "templates" for actual camera instances: Whenever a node refers to one
> camera, a new instance of this camera is created.
>
> There is no "default" camera for a glTF asset.
> Instead, client application has to keep track of the currently active camera.
> The client application may, for example, offer a dropdown-menu that allows
> one to select the active camera and thus to quickly switch between predefined
> view configurations. With a bit more implementation effort, the client
> application can also define its own camera and interaction patterns for the
> camera control (e.g., zooming with the mouse wheel).

Interesting that camera is not specified, it is up to client. I think I should
adopt this format as well.

This will be a modest change to my scene format.
That means my camera node? So and lighting too?

I've started implementing my glTF package, but until I really start loading
files and messing with it, I don't think I will understand it from all angles.


# 2024-11-20 faster than the speed of light

Inspired by the success of getting texture mapping working, I focused on some
key performance improvements. It was clear from profiling that my BVH partition
strategy was suboptimal. I completely overhauled this function with a very
comprehensive binned SAH implementation. This singlehandedly resulted in over
a 2x speedup in rendering time. For the time being, at least, I found that
the trace profiling was not as useful as I had expected. Instead, the CPU
profiling is leading me directly to the source of the slowdown.

In total, I managed to bring render time down from 12us per pixel to now just
2.6us per pixel.

I am really stoked that I finally have a working render.

Render belows a camera orbiting a 3D scan model of a spray bottle, varying
in angle and distance.

<img src="./scene\scan\out\out_20241119_161650.gif">

What's awesome is that it no longer consumes ridiculous RAM like it did with
tracing. I think for my application, a purpose built ultra light logger is the
way to go.

Really satisfied with the refactor that I did today. I made an enormous amount
of changes throughout the codebase, and now with my significantly improved BVH
implementation, I am able to render 2048x2048 images in just a few seconds,
whereas before it would take the same time to render a 512x512 image.

I temporarily broke animations, however, and I'll have to go back and fix my old scene files to use whatever new solution I come up with.


# 2024-11-22 adding an interactive gui to my phys renderer
I'm now starting to port some code files from batleaxe over to phys.
I created a new package "lab" as a the interactive playground for the phys
renderer package.

After porting the initial files, I'm thinking about my next steps and my vision
for the interactive gui in general.

Inspired ny the recent glTF talks, I like the concept an "asset" as bundling
related scenes and multiple cameras, instead of just a single scene and camera.

Then it is up to the viewer to decide which camera to use and which scene to
render, although there is a default scene field in the asset. I also like that
when there are no scenes, then the gltf asset is interpreted as a library
of entities like materials or meshes. That's really clever and elegant.

What I like about this representation is that it holds enough information for
it to really start to feel like a standalone asset.

Even if, at least initially, I don't export or import glTF files, I think the
conceptual model is worth adopting. For example, I could still pretend that
an OBJ file and accompanying MTL file are a glTF asset.

So what the heck do I want to do exactly? How can I break the next steps down
into smaller tasks? That's the part I'm stuck on at the moment.

I remember the booklet "how to solve it" by George Polya. What's my goal?

Everything is easier when I think of phys/lab as a viewer, at least for now.
Then I can think of the viewer as a standalone application that can load
assets and render them.

Extending from there, I can display the scene graph somewhere on the page, not
necessarily the canvas, but as tabs or a tree view.

One advantage to keeping my scene small and flat is that I can easily map it
to a table view.

So what do I even need the websocket server for?
I think I can put the entire renderer on the client side, so my server is
responsible for serving static assets including the client wasm binary.

Hmm it is becoming more clear how the responsibilities are divided now.

The client is js wasm binary. The client refers to assets by url and assumes
there exists a server that will resolve them. The client is responsible for
setting up screen interface and rendering the scene.
The client initially loads no model, but can make a GET to the server to query
it for available scenes and assets. Then the client when it picks what it
wants to show then sends GET requests to the server to load the assets, which
may be individual asset textures or a larger conceptual object like a gltf,
containing multiple scenes and cameras.

The client imports slam/code/photon/raytrace/phys and uses it to render the scene.
The client also implements essential camera controls like orbiting and zooming.
For example, dragging to orbit will continuously update the camera position
and then call the render function. It is great that the renderer supports
cancellation so well, since I may want to interrupt the status of the render
at any time.

While rendering, the tiles stream out as completed events, which fills out the
canvas that the client is rendering to.

When the client is handling frequent events such as being dragged continuously,
the client can throttle the number of render calls to avoid overloading, and
the client may temporarily reduce the resolution of the render to improve
performance. After the client stops dragging, the client can then request a
final render at full resolution. When this happens, the resolution of the canvas
actually never changes, only the image that is mapped to it. The canvas is
always set to render resolution, but when temporarily reducing resolution,
the client can render to an offscreen canvas and then scale it down before
mapping it to the main canvas.

This means one of the next steps will be to add this gltf style asset type to
my phys package.

Sketching out the most critical functionality:

0. script builds wasm binary and prepares the dist folder
1. server launches, serves static files from dist folder, and serves an asset directory listing
2. client loads wasm binary, queries server for available assets, and renders them

All of the other things like orbit interaction and throttling can wait until
the critical path has been implemented.

The critical tasks are listed below for the client and server:

Server:

1. listenAndServe on $PORT or default to 8060
2. handle GET "/" with file server from the dist folder.
3. handle GET "/assets" with a list of available assets in the dist folder.

no websocket server needed, and the file server from dist handles all of the
assets as well, very simple!.

Client:

1. connect to canvas and wasm binary
2. query server for available assets
3. render the first asset


This is a good first step, if I can accomplish this, I think I can iterate
on it to add the remaining features.

That means I'll gut the entire server package and use main.go as the server
top level main which hosts the binary. Really makes things quite a bit easier.
Websockets would be really neat in a sense because it would make live collab
possible, but I think that's a feature that can be added later.

# 2024-11-22 more cool finds in the go mobile standard library

In `go/x/mobile` there are lots of really interesting packages for problems
related to packaging and shipping go code to mobile devices, but more generally
any platform using the device interface.

There are some of my favorite finds:

`go/x/mobile/asset` - Package for loading named assets. This would be useful to
me in the renderer, battleaxe, and lab frontend. I like that it uses build
constraints to select the appropriate implementation for the platform.
I'm not sure that I will use build constraints for me, but defining an asset
interface may be worthwhile at some point. Let's see what problems actually
show up as I develop it more.
https://cs.opensource.google/go/x/mobile/+/fa514ef7:asset/asset.go

In particular, I really like the pump function in `go/x/mobile/app`.
When you get into the details of channels in the context of apps like
multiplayer servers, or rendering, the pump function is a very convenient way
to send events around the app without worrying about blocking or deadlocks.
I really like the function and the comment docstring so I'm including it here.

I've run into exactly this problem before in `slam/code/battleaxe` and I had
a similar function defined but not as clean or as coherent as the focus in this
one.

What I like is that it splits eventsIn and eventsOut so that we can work with
the channel most effectively. That's something I struggled with when writing
the websocket netcode for battleaxe.

```go
type stopPumping struct{}

// pump returns a channel src such that sending on src will eventually send on
// dst, in order, but that src will always be ready to send/receive soon, even
// if dst currently isn't. It is effectively an infinitely buffered channel.
//
// In particular, goroutine A sending on src will not deadlock even if goroutine
// B that's responsible for receiving on dst is currently blocked trying to
// send to A on a separate channel.
//
// Send a stopPumping on the src channel to close the dst channel after all queued
// events are sent on dst. After that, other goroutines can still send to src,
// so that such sends won't block forever, but such events will be ignored.
func pump(dst chan interface{}) (src chan interface{}) {
	src = make(chan interface{})
	go func() {
		// initialSize is the initial size of the circular buffer. It must be a
		// power of 2.
		const initialSize = 16
		i, j, buf, mask := 0, 0, make([]interface{}, initialSize), initialSize-1

		srcActive := true
		for {
			maybeDst := dst
			if i == j {
				maybeDst = nil
			}
			if maybeDst == nil && !srcActive {
				// Pump is stopped and empty.
				break
			}

			select {
			case maybeDst <- buf[i&mask]:
				buf[i&mask] = nil
				i++

			case e := <-src:
				if _, ok := e.(stopPumping); ok {
					srcActive = false
					continue
				}

				if !srcActive {
					continue
				}

				// Allocate a bigger buffer if necessary.
				if i+len(buf) == j {
					b := make([]interface{}, 2*len(buf))
					n := copy(b, buf[j&mask:])
					copy(b[n:], buf[:j&mask])
					i, j = 0, len(buf)
					buf, mask = b, len(b)-1
				}

				buf[j&mask] = e
				j++
			}
		}

		close(dst)
		// Block forever.
		for range src {
		}
	}()
	return src
}
```

`go/x/mobile/app/darwin_desktop.go` - Example of making



# 2024-11-23 rendering loop on the phys lab client

I now have a basic client that can render a triangle.
The next task is to implement a rendering loop that can render a scene.
There are multiple different components to implement here.
I need a way to view and edit scene graphs.
I need a way to see the render output.
I need a way to move the camera with the mouse (orbit style controls, zoom, pan).
When either the scene graph or camera changes, I need to trigger a render.

I like the go/x/mobile/app structure, and I would benefit from collecting all
of my globals and events into some easy to analyze streams.



# 2024-11-23 refresher on the coordinate systems in js wasm

Viewport
1. The viewport is the rectangle on the screen where the image is drawn.
2. MouseEvent.clientX and MouseEvent.clientY
3. TouchEvent.clientX and TouchEvent.clientY

Page
1. The page is the entire document. Increases as you scroll down or right.
2. MouseEvent.pageX and MouseEvent.pageY
3. TouchEvent.pageX and TouchEvent.pageY

Screen
1. The screen is the entire screen. It is the same as the window size.
2. MouseEvent.screenX and MouseEvent.screenY
3. TouchEvent.screenX and TouchEvent.screenY


# 2024-11-25 a physical button for recording and transcribing my voice

As I use AI tools more and more, I'm starting to think about how I can use
them even faster and more effectively. One of the annoying things sometimes
is that it takes a while to write a prompt. Sometimes I am limited by the
speed of my thought, and sometimes I'm limited by the speed that I can type,
which to be honest, is quite fast.

I'm thinking about how I might set up a physical button that I can press to
quickly record a voice memo and then automatically transcribe it to text using
a model such as OpenAI whisper. Another use case is recording my thoughts in
text such as this entry. I anticipate that if I had a solution like this, I
would be able to write many more tokens of text and document my thoughts more
for my personal projects.

It has to be fast and easy to use, and it has to automatically convert the
speech to text quickly and accurately so that I can paste it into my prompts.

To avoid having to put on and off my headphones, and take several seconds to
turn them on and off, I might want to have an additional dedicated microphone.

Today I installed the openai whisper model and it was surprisingly easy to
get started with.

To install the model, I ran:

```
py -m pip install git+https://github.com/openai/whisper.git
```

I already had `ffmpeg` installed, which is required to use whisper.

After installing, I have two ways of running the model.
The first way is with the new whisper command line tool.

```
whisper audio.flac audio.mp3 audio.wav --model turbo
```

The second is with the python module:

```
import whisper

model = whisper.load_model("turbo")
result = model.transcribe("audio.mp3")
print(result["text"])
```


# 2024-11-25 thinking about next steps for my renderer

I'm thinking about where to go next with my renderer. I have a basic renderer
working now, it does 3D, it has panning and zooming controls. I can load .obj
and .mtl files with textures. I can render scenes with multiple objects.

I'm a bit disappointed with the speed of the interactive GUI, mainly because
the WASM is limited to a single thread. That means I'm only getting one CPU to
max out, whereas when I run it on desktop I can fully utilize all of my cores.

I don't have volume rendering and I don't have anything except for lambertian
shading.

I'll break down my thoughts for improvements by package.

client:

1. moving the camera should interrupt the current render and start a new one.
2. rendering should stream tiles to the canvas instead of drawing at the end.
3. toolbar with toggle for axes arrows
4. cancelation all the way down
5. loading screen gone only after assets loaded
6. events don't block
7. scene/camera picker
8. scene graph viewer
9. web worker to render in background and use multiple cpu cores
10. scene selector to change what is displayed
11. binary encoding with worker

phys:

1. gltf support
2. gpu rendering
3. volume rendering
4. quaternion rotations
5. ray profiling
6. clean cancelation all the way down
7. better progress reporting
8. log errors and render events
9. pbr rendering with metalness and roughness
10. binary encoding


things to ask for help on:

1. how to adapt my code for gpu rendering


# 2024-11-26 more thoughts about next steps for learning and renderer

I'm at a bit of a crossroads with my renderer. I have achieved all of my
original goals, I have textures, 3D rendering, direct lights, and a basic
interactive GUI. In the immediate short term, I'll continue to work on
refactoring and cleaning up the codebase, but I'm still thinking about my next
big milestone to work towards. I've thought about a range of different things,
such as a basic FPS kind of game, learning GPU rendering, or glTF support.

At the moment, I'm leaning towards learning more about GPU rendering or glTF.
To be honest, because of how many overlapping concepts there are in both, I
don't think it really matters which one I pick, they will both be useful to
me in the long term, and I will likely end up learning both.

Today I refactored my frontend code to use a web worker, which was quite a
journey to make it work, and in the end, I still have more to do.

The biggest issue was the message passing. Because my scene is detailed,
converting to JSON takes a long time, like nearly 10 minutes in one case.
I think the solution is to put the core of the app in the worker, and in the
main thread I will do event handling and send render start events.


# 2024-11-26 scott's webgl2 notes

Documenting some of my webgl2 notes here.

Good resource for webgl: https://webgl2fundamentals.org/
Since 2021, all major browsers support webgl2.

webgl2 state diagram (really cool!):
https://webgl2fundamentals.org/webgl/lessons/resources/webgl-state-diagram.html

1. programs reference textures by index.
2. buffers and vertex arrays attribute are also referenced by index.
3. webgpu is successor but has partial support

# 2024-11-28 camera drivers for lab package

As I think less now about working on the renderer, and
shifting a bit more to focus on fourier imaging methods,
in particular, particle sizing using a method similar to
differential dynamic microscopy.

In the past, I've focused less on the speed of my camera driver code and
more on the correctness and declarative interface. The particule sizing
applicaton is interesting because it is the first time I've encountered a
situation that calls for very rapid continuous capture of the same camera
settings. For most people, this is the typical way to use a camera.
In this particular situation, sensor spatial resolution is not as important
as framerate.

I realized that my existing camera driver code doesn't have a solution for
continuous capture yet, so today I am revisiting the camera driver code to
make another round of improvements, and also sketch out a design for a camera
interface that works for my personal use cases.

While going back and forth on different ideas for Go Camera interfaces, I
wrote down some mostly random and incomplete rambling thoughts about the
camera interface which I've included below for reference.

```
Thank you for your thoughtful answer. I like the Camera interface for the most part but I have some thoughts. One is that pointers can be really annoying to work with for the Config struct, but I think we could carefully define the values such that the default zero value means "unset". That works elegantly for exposure, since a valid exposure must be nonzero, it works for gain, assuming a certain interpretation of gain where it means the multiplicative factor like 1x or 2x. BlackLevel it also works because zero black level means, well, no black level! And the zero value ROI struct also can be seen as unset.

I've decided to rename StartAcquisition to simply Start, and likewise with StopAcquisition to Stop.

For now, I've decided to not include triggering at all from the interface, and don't include it. So that just leaves Configure, Capture, Start, Stop.

I've decided to omit the cfg Config from Capture, thereby making capture's interpretation that it will always capture exactly one snapshot of whatever the current camera configuration is.

I thought it was interesting how you created a Frame struct with a timestamp and raw byte buffer. I like that because it is low level enough to cover most cameras. I also thought it was interesting how you added a metadata field, I can imagine that being really useful.

I now want to ask you about some other aspects. First is that since this is a camera, it is a high performance device and I want to be able to understand how fast the main camera operations are. That means I want to know how long it takes to set a camera node, how long it takes to set exposure, and how long it takes to readout the frame data, and how long it takes to convert to my eventual final image format. Would you revise anything about the camera interface with this in mind?

Second aspect is testing. How can I test implementations of the camera interface, and how can I make a fake camera so that I can develop camera based applications without needing a physical hardware camera connected?

Third aspect is opening and closing the camera connection. I look at this from many perspectives. From one angle, it is nice to leave the details of opening the connection to the setup code, and then focus on the camera interface as only the part where we actually use it. But another perspective would be to add a device open and close function to the interface so that the user can cleanup the interface properly when done. Another perspective would be to be like a ReadCloser and only include the close method in the interface, but not open. Please comment on some good options to consider here.

Fourth aspect is device events and errors. Imagine we call Start and then read out a few frames, suddenly the device is unplugged, now something will panic or raise and error, but how is that reflected and defined in the context of our interface?

Fifth aspect is about metadata, what do you think I might put here? I find it really interesting as an idea. How would I serialize this when I save images? Does using metadata influence what image formats I can use, based on their respective support for metadata?

Sixth aspect is about rethinking triggers. I want to think about triggers not as part of the general camera interface, but as part of only a particular SDK implementation, at least for now. But even then, I think about whether triggers is really modelled well here. What we are doing with triggers is setting up source line and source activation condition. Then we wait for that event to happen and trigger the acquisition start event. Then wait until the acquisition finish event and grab the image.

In your response, please carefully consider the thoughtful aspects that I have pointed out. Do not feel obligated to follow or agree with me, I am looking for a thoughtful, insightful, careful, and critical examination.

Remember that the overall goal is to consider my use case and what I'm trying to model and use these cameras for, then to infer what aspects are best included in the go interface, and to think though the event and behaviours that are expected. In your response, please carefully consider my points and give your thoughtful response.
```

This really discusses a Camera interface that looks like this:

```go

type Config struct {
  Exposure time.Duration
  Gain     float64
  ROI      image.Rectangle
  Color    color.Model
}

type Camera interface {
    Configure(cfg Config) error
    Capture() (Frame, error)
    Start(ctx context.Context) (<-chan Frame, error)
    Stop() error
}
```

