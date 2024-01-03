#!/usr/bin/env python
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP

import os
import sys


def main(args):
    n = len(args)
    if n != 2:
        print("Pass the log directory")
        return 1
    log_path = args[1]
    file_content = ''
    for x in os.listdir(log_path):
        if x.endswith(".txt"):            
            with open(os.path.join(log_path, x)) as f:
                file_content += (f.read())
    test_count = file_content.count('RUN')
    pass_count = file_content.count('PASS:')
    fail_count = file_content.count('FAIL:')
    print(
        f"\nTestcases Ran: {test_count}; "
        f"Testcases Passed: {pass_count}; "
        f"Testcases Failed: {fail_count}; ")
    return 0

if __name__ == "__main__":
    exit(main(sys.argv))
