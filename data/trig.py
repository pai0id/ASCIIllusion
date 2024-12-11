def triangulate_faces(obj_file, output_file):
    vertices = []
    normals = []
    texture_coords = []
    faces = []

    # Read the OBJ file
    with open(obj_file, 'r') as file:
        for line in file:
            if line.startswith('v '):
                # Vertex
                vertices.append(line.strip())
            elif line.startswith('vn '):
                # Vertex normal
                normals.append(line.strip())
            elif line.startswith('vt '):
                # Texture coordinate
                texture_coords.append(line.strip())
            elif line.startswith('f '):
                # Face
                face = line.strip().split()[1:]
                faces.append(face)

    # Triangulate faces
    triangulated_faces = []
    for face in faces:
        if len(face) > 3:
            # Polygon face, split into triangles
            for i in range(1, len(face) - 1):
                triangulated_faces.append([face[0], face[i], face[i + 1]])
        else:
            # Already a triangle
            triangulated_faces.append(face)

    # Write the triangulated OBJ file
    with open(output_file, 'w') as file:
        # Write vertices
        for vertex in vertices:
            file.write(f"{vertex}\n")

        # Write normals
        for normal in normals:
            file.write(f"{normal}\n")

        # Write texture coordinates
        for texture in texture_coords:
            file.write(f"{texture}\n")

        # Write faces
        for face in triangulated_faces:
            file.write(f"f {' '.join(face)}\n")

# Example usage
input_obj = "dingus.obj"
output_obj = "dinguser.obj"
triangulate_faces(input_obj, output_obj)

