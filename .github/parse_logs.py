#!/usr/bin/env python
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP

import os
import sys


def main(args):
    n = len(args)
    if n != 2:
        print("Pass the log directory or txt file path")
        return 1
    log_path = args[1]
    file_content = ''
    if log_path.endswith(".txt"):
        with open(log_path) as f:
            file_content += (f.read())
    else:
        for x in os.listdir(log_path):
            if x.endswith(".txt"):            
                with open(os.path.join(log_path, x)) as f:
                    file_content += (f.read())
    test_count = file_content.count('RUN') - file_content.count('SKIP:')
    pass_count = file_content.count('PASS:')
    fail_count = file_content.count('FAIL:')
    print(
        f"\nTestcases Ran: {test_count}; \n"
        f"Testcases Passed: {pass_count}; \n"
        f"Testcases Failed: {fail_count}; \n")
    return 0

if __name__ == "__main__":
    exit(main(sys.argv))
