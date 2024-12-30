import os
import re
import tempfile
import shutil

def remove_normals_and_textures_from_obj(input_file, output_file):
    try:
        with open(input_file, 'r') as infile, \
             (tempfile.NamedTemporaryFile('w', delete=False) if input_file == output_file else open(output_file, 'w')) as outfile:
            temp_name = outfile.name
            for line in infile:
                if line.startswith('vn ') or line.startswith('vt '):
                    continue
                if line.startswith('f '):
                    cleaned_line = re.sub(r'(\d+)/\d*/\d*', r'\1', line)
                    outfile.write(cleaned_line)
                else:
                    outfile.write(line)

        if input_file == output_file:
            try:
                os.replace(temp_name, input_file)
            except OSError as e:
                if e.errno == 18:  # Invalid cross-device link
                    shutil.move(temp_name, input_file)
                else:
                    raise

        print(f"Normals and texture coordinates removed. Output saved to {output_file}")
    except FileNotFoundError:
        print(f"Error: File '{input_file}' not found.")
    except IOError as e:
        print(f"An I/O error occurred: {e}")

def triangulate_obj(input_file, output_file):
    try:
        with open(input_file, 'r') as infile, \
             (tempfile.NamedTemporaryFile('w', delete=False) if input_file == output_file else open(output_file, 'w')) as outfile:
            temp_name = outfile.name
            for line in infile:
                if line.startswith('f '):
                    parts = line.split()
                    face_vertices = parts[1:]

                    if len(face_vertices) <= 3:
                        outfile.write(line)
                    else:
                        for i in range(1, len(face_vertices) - 1):
                            triangle = f"f {face_vertices[0]} {face_vertices[i]} {face_vertices[i + 1]}\n"
                            outfile.write(triangle)
                else:
                    outfile.write(line)

        if input_file == output_file:
            try:
                os.replace(temp_name, input_file)
            except OSError as e:
                if e.errno == 18:  # Invalid cross-device link
                    shutil.move(temp_name, input_file)
                else:
                    raise

        print(f"Successfully triangulated and saved to {output_file}")
    except FileNotFoundError:
        print(f"Error: The file '{input_file}' does not exist.")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

def main():
    print("Welcome to the OBJ Cleaner Tool!")
    input_path = input("Enter the path to the .obj file: ").strip()
    output_path = input("Enter the path to save the output file: ").strip()
    triangulate = input("Triangulate? (Y/n): ").strip()
    
    if not os.path.exists(input_path):
        print(f"Error: The file '{input_path}' does not exist.")
        return
    
    remove_normals_and_textures_from_obj(input_path, output_path)

    if triangulate.lower() != 'n':
        triangulate_obj(output_path, output_path)

if __name__ == "__main__":
    main()
