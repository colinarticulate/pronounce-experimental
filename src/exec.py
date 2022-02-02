'''
    Module to execute commands on the shell. Command line is send as a list of one single string as it would have been typed in the command line.
'''

import subprocess


def execute(command, cwd, file=None):

    if file != None and file != "":
        with open(file, 'w') as f:
            process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, cwd=cwd, universal_newlines=True)

            while True:
                output = process.stdout.readline()
                print(output.strip())
                f.write(output)
                # Do something else
                return_code = process.poll()
                if return_code is not None:
                    print('\n>>> RETURN CODE', return_code)
                    f.write(f"\n>>> RETURN CODE {return_code}\n")
                    # Process has finished, read rest of the output 
                    for output in process.stdout.readlines():
                        print(output.strip())
                        f.write(output)
                    break

    else:
        process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, cwd=cwd, universal_newlines=True)

        while True:
            output = process.stdout.readline()
            print(output.strip())

            # Do something else
            return_code = process.poll()
            if return_code is not None:
                
                print('\n>>> RETURN CODE', return_code)

                # Process has finished, read rest of the output 
                for output in process.stdout.readlines():
                    print(output.strip())

                break
