#!/usr/bin/env python3
import subprocess
from ruamel.yaml import YAML
import glob
import re
import argparse

class Semver:
    def __init__(self, major: int, minor: int, patch: int):
        self.major = major
        self.minor = minor
        self.patch = patch

    def incr_major(self):
        self.major = self.major + 1
        self.patch = 0
        self.minor = 0

    def incr_minor(self):
        self.minor = self.minor + 1
        self.patch = 0

    def incr_patch(self):
        self.patch = self.patch + 1

    def to_string(self) -> str:
        return F'{self.major}.{self.minor}.{self.patch}'

    @classmethod
    def parse(self, input_str: str):
        # Parse and validate, return new instance of Semver
        if re.fullmatch(r'^([0-9]+)\.([0-9]+)\.([0-9]+)$', input_str):
            split_list = input_str.split('.')
            split_list = [int(i) for i in split_list]
            return self(*split_list)

        raise Exception(F'Invalid input version value: {input_str}')


def update_charts(update_func):
    yaml = YAML()
    yaml.indent(mapping = 2, sequence=4, offset=2)
    main_dir = subprocess.run(["git", "rev-parse", "--show-toplevel"], check=True, stdout=subprocess.PIPE).stdout.strip().decode('utf-8')
    search_path = F'{main_dir}/charts/**/Chart.yaml'
    for path in glob.glob(search_path, recursive=True):

        if re.match('^.*cass-operator.*', path):
            continue

        with open(path) as f:
            chart = yaml.load(f)

        semver = Semver.parse(chart['version'])
        update_func(semver)
        chart['version'] = semver.to_string()        

        with open(path, 'w') as f:
            yaml.dump(chart, f)

        print(F'Updated {path} to {semver.to_string()}')

def main():
    parser = argparse.ArgumentParser(description='Update Helm chart versions in k8ssandra project')
    parser.add_argument('--incr', choices=['major', 'minor', 'patch'], help='increase part of semver by one')
    args = parser.parse_args()

    if args.incr:
        if args.incr == 'major':
            update_charts(Semver.incr_major)
        elif args.incr == 'minor':
            update_charts(Semver.incr_minor)
        elif args.incr == 'patch':
            update_charts(Semver.incr_patch)

if __name__ == "__main__":
    main()
