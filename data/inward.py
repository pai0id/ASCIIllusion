import numpy as np

def parse_obj(file_path):
    """Parse vertices and normals from an OBJ file."""
    vertices = []
    normals = []
    faces = []
    with open(file_path, 'r') as file:
        for line in file:
            parts = line.split()
            if not parts:
                continue
            if parts[0] == 'v':  # Vertex
                vertices.append([float(x) for x in parts[1:]])
            elif parts[0] == 'vn':  # Normal
                normals.append([float(x) for x in parts[1:]])
            elif parts[0] == 'f':  # Face
                faces.append(parts[1:])
    return np.array(vertices), np.array(normals), faces

def calculate_centroid(vertices):
    """Calculate the centroid of the object."""
    return np.mean(vertices, axis=0)

def adjust_normals(vertices, normals):
    """Adjust normals to face inward."""
    centroid = calculate_centroid(vertices)
    adjusted_normals = []
    for vertex, normal in zip(vertices, normals):
        direction_to_centroid = centroid - vertex
        if np.dot(normal, direction_to_centroid) > 0:
            adjusted_normals.append(-normal)  # Reverse the normal
        else:
            adjusted_normals.append(normal)
    return np.array(adjusted_normals)

def write_obj(file_path, vertices, normals, faces):
    """Write the adjusted OBJ file."""
    with open(file_path, 'w') as file:
        for vertex in vertices:
            file.write(f"v {' '.join(map(str, vertex))}\n")
        for normal in normals:
            file.write(f"vn {' '.join(map(str, normal))}\n")
        for face in faces:
            file.write(f"f {' '.join(face)}\n")

def main(input_path, output_path):
    vertices, normals, faces = parse_obj(input_path)
    adjusted_normals = adjust_normals(vertices, normals)
    write_obj(output_path, vertices, adjusted_normals, faces)
    print(f"Adjusted OBJ file saved to {output_path}")

if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description="Adjust OBJ normals to face inward.")
    parser.add_argument("input", help="Path to the input OBJ file.")
    parser.add_argument("output", help="Path to save the adjusted OBJ file.")
    args = parser.parse_args()
    main(args.input, args.output)
