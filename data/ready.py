import sys
import argparse

def read_obj(file_path):
    vertices = []
    normals = []
    faces = []
    
    with open(file_path, 'r') as file:
        for line in file:
            if line.startswith('v '):
                parts = line.strip().split()
                vertices.append(tuple(map(float, parts[1:])))
            elif line.startswith('vn '):
                parts = line.strip().split()
                normals.append(tuple(map(float, parts[1:])))
            elif line.startswith('f '):
                parts = line.strip().split()
                face = []
                for p in parts[1:]:
                    indices = p.split('/')
                    vertex_idx = int(indices[0])
                    normal_idx = int(indices[2]) if len(indices) > 2 and indices[2] else None
                    face.append((vertex_idx, normal_idx))
                faces.append(face)
    
    return vertices, normals, faces

def triangulate_faces(faces):
    triangulated_faces = []

    for face in faces:
        if len(face) == 3:
            triangulated_faces.append(face)
        else:
            for i in range(1, len(face) - 1):
                triangulated_faces.append([face[0], face[i], face[i + 1]])

    return triangulated_faces

def write_obj(file_path, vertices, normals, faces):
    with open(file_path, 'w') as file:
        for vertex in vertices:
            file.write(f"v {vertex[0]} {vertex[1]} {vertex[2]}\n")

        for normal in normals:
            file.write(f"vn {normal[0]} {normal[1]} {normal[2]}\n")

        for face in faces:
            face_str = ' '.join(
                f"{v_idx}//{n_idx}" if n_idx else f"{v_idx}"
                for v_idx, n_idx in face
            )
            file.write(f"f {face_str}\n")

def main():
    parser = argparse.ArgumentParser(description="Triangulate faces in an OBJ file, including normals.")
    parser.add_argument("input_file", help="Path to the input OBJ file.")
    parser.add_argument("output_file", help="Path to the output OBJ file.")
    args = parser.parse_args()

    input_file = args.input_file
    output_file = args.output_file

    try:
        vertices, normals, faces = read_obj(input_file)
        triangulated_faces = triangulate_faces(faces)
        write_obj(output_file, vertices, normals, triangulated_faces)
        print(f"Triangulated OBJ file with normals saved to {output_file}")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()