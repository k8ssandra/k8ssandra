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

def update_chart_version(update_func, chart_name, update_dep):
        yaml = YAML()
        yaml.indent(mapping = 2, sequence=4, offset=2)
        main_dir = subprocess.run(["git", "rev-parse", "--show-toplevel"], check=True, stdout=subprocess.PIPE).stdout.strip().decode('utf-8')
        path = F'{main_dir}/charts/{chart_name}/Chart.yaml'
        with open(path) as f:
            chart = yaml.load(f)

        semver = Semver.parse(chart['version'])
        update_func(semver)
        chart['version'] = semver.to_string()        

        with open(path, 'w') as f:
            yaml.dump(chart, f)

        print(F'Updated {path} to {semver.to_string()}')
        if chart_name != 'k8ssandra' and update_dep is True:
            update_dependency(main_dir, chart_name, semver.to_string())
            update_chart_version(update_func, 'k8ssandra', False)

def update_dependency(main_dir, dependency, target_version):
    yaml = YAML()
    yaml.indent(mapping = 2, sequence=4, offset=2)
    path = F'{main_dir}/charts/k8ssandra/Chart.yaml'
    # If we update the k8ssandra dependency.. we need to update its version also. Why not do it with the same code as above?
    with open(path) as f:
        chart = yaml.load(f)

    for k in chart['dependencies']:
        if k['name'] == dependency:
            k['version'] = target_version
            with open(path, 'w') as f:
                yaml.dump(chart, f)

            print(F'Updated dependency {dependency} to {target_version}')
            break

def main():
    parser = argparse.ArgumentParser(description='Update Helm chart versions in k8ssandra project')
    parser.add_argument('--incr', choices=['major', 'minor', 'patch'], help='increase part of semver by one')
    parser.add_argument('--chart', choices=['medusa-operator','k8ssandra-common','reaper-operator','k8ssandra','cass-operator','backup','restore'], help='increase only a single chart')
    parser.add_argument('--dep', help='update the k8ssandra dependency to newest version also', dest='dep_update', action='store_true')
    args = parser.parse_args()

    chart = args.chart
    if not args.chart:
        chart = "*"

    update_dep = False

    if args.dep_update is True:
        update_dep = True

    if args.incr:
        if args.incr == 'major':
            update_chart_version(Semver.incr_major, chart, update_dep)
        elif args.incr == 'minor':
            update_chart_version(Semver.incr_minor, chart, update_dep)
        elif args.incr == 'patch':
            update_chart_version(Semver.incr_patch, chart, update_dep)

if __name__ == "__main__":
    main()
