import sys
import os
import math

class Vec3:
    def __init__(self, x, y, z):
        self.x = x
        self.y = y
        self.z = z

    def normalize(self):
        length = math.sqrt(self.x**2 + self.y**2 + self.z**2)
        if length > 0:
            self.x /= length
            self.y /= length
            self.z /= length

    def negate(self):
        self.x = -self.x
        self.y = -self.y
        self.z = -self.z

    def __str__(self):
        return f"{self.x:.6f} {self.y:.6f} {self.z:.6f}"

class Face:
    def __init__(self, vertices):
        self.vertices = vertices  # List of tuples: (vertex_index, texture_index, normal_index)

    def triangulate(self):
        if len(self.vertices) <= 3:
            return [self]
        triangles = []
        for i in range(1, len(self.vertices) - 1):
            triangles.append(Face([self.vertices[0], self.vertices[i], self.vertices[i + 1]]))
        return triangles

    def __str__(self):
        return "f " + " ".join(
            ["/".join(map(lambda x: str(x) if x is not None else "", vertex)) for vertex in self.vertices]
        )

def parse_obj(file_path):
    vertices = []
    normals = []
    faces = []

    with open(file_path, 'r') as file:
        for line in file:
            parts = line.strip().split()
            if not parts or parts[0].startswith('#'):
                continue
            if parts[0] == 'v':
                vertices.append(Vec3(*map(float, parts[1:])))
            elif parts[0] == 'vn':
                normals.append(Vec3(*map(float, parts[1:])))
            elif parts[0] == 'f':
                face_vertices = []
                for part in parts[1:]:
                    indices = part.split('/')
                    vertex = tuple(int(indices[i]) if i < len(indices) and indices[i] else None for i in range(3))
                    face_vertices.append(vertex)
                faces.append(Face(face_vertices))

    return vertices, normals, faces

def write_obj(file_path, vertices, normals, faces):
    with open(file_path, 'w') as file:
        for v in vertices:
            file.write(f"v {v}\n")
        for n in normals:
            file.write(f"vn {n}\n")
        for f in faces:
            file.write(f"{f}\n")

def calculate_normal(face, vertices):
    v1 = vertices[face.vertices[0][0] - 1]
    v2 = vertices[face.vertices[1][0] - 1]
    v3 = vertices[face.vertices[2][0] - 1]

    ux, uy, uz = v2.x - v1.x, v2.y - v1.y, v2.z - v1.z
    vx, vy, vz = v3.x - v1.x, v3.y - v1.y, v3.z - v1.z

    nx, ny, nz = uy * vz - uz * vy, uz * vx - ux * vz, ux * vy - uy * vx
    normal = Vec3(nx, ny, nz)
    normal.normalize()
    normal.negate()
    return normal

def process_obj(input_file, output_file):
    vertices, normals, faces = parse_obj(input_file)

    if not normals:
        normals = []
        for face in faces:
            for vertex in face.vertices:
                normal = calculate_normal(face, vertices)
                normals.append(normal)

    new_faces = []
    for face in faces:
        triangles = face.triangulate()
        for triangle in triangles:
            new_faces.append(triangle)

    write_obj(output_file, vertices, normals, new_faces)

def main():
    if len(sys.argv) != 3:
        print("Usage: python process_obj.py <input_file> <output_file>")
        sys.exit(1)

    input_file = sys.argv[1]
    output_file = sys.argv[2]

    if not os.path.exists(input_file):
        print(f"Error: Input file '{input_file}' does not exist.")
        sys.exit(1)

    process_obj(input_file, output_file)
    print(f"Processed OBJ saved to '{output_file}'")

if __name__ == "__main__":
    main()
