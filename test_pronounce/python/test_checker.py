import os

test_results_dir = "../output_testing_xyz_plus"
summary_file = os.path.join(test_results_dir,"000__summary__000.txt")

files_dir = [f.split(".")[0] for f in sorted(os.listdir(test_results_dir))[1:]]

with open(summary_file,'r') as f:
    contents = f.read()
files_summary = [f.split(",")[0] for f in contents.strip("\n").split("\n")[:-1] ]

missing=[]
for file in files_dir:
    if file not in files_summary:
        missing.append(file)
print("Working on it")