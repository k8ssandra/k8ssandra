#!/usr/bin/env python3
import subprocess
from ruamel.yaml import YAML
import argparse


def update_tag_version(chart, component, tag):
    yaml = YAML()
    yaml.indent(mapping=2, sequence=4, offset=2)
    main_dir = subprocess.run(["git", "rev-parse", "--show-toplevel"],
                              check=True, stdout=subprocess.PIPE).stdout.strip().decode('utf-8')
    path = F'{main_dir}/charts/{chart}/values.yaml'
    with open(path) as f:
        values = yaml.load(f)

    if component is not None and component is not "":
        values[component]['image']['tag'] = tag
    else:
        values['image']['tag'] = tag

    with open(path, 'w') as f:
        yaml.dump(values, f)

    if component is not None and component is not "":
        print(F'Updated {chart}/{component} to use {tag}')
    else:
        print(F'Updated {chart} to use {tag}')


def main():
    parser = argparse.ArgumentParser(
        description='Update image tag version in k8ssandra project')
    parser.add_argument('--tag', dest='image_tag', action='store')
    parser.add_argument('--component', dest='component', action='store')
    parser.add_argument('--chart', choices=['k8ssandra', 'medusa-operator', 'reaper-operator',
                                            'cass-operator'], help='target operator image version to update')
    args = parser.parse_args()

    if args.image_tag and args.chart:
        update_tag_version(args.chart, args.component, args.image_tag)


if __name__ == "__main__":
    main()
