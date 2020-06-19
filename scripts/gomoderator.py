import os
import sys
import subprocess
import errno
from typing import List, Dict

"""
Resolves dependency conflicts between a plugin repository's and the core repository's go.mods

Usage: python3 gomoderator.py {path_to_core_repository} {path_to_plugin_repository}
"""

ERROR_INVALID_NAME = 123


def is_pathname_valid(pathname: str) -> bool:
    """
    `True` if the passed pathname is a valid pathname for the current OS;
    `False` otherwise.
    """
    try:
        if not isinstance(pathname, str) or not pathname:
            return False
        _, pathname = os.path.splitdrive(pathname)
        root_dirname = os.environ.get('HOMEDRIVE', 'C:') \
            if sys.platform == 'win32' else os.path.sep
        assert os.path.isdir(root_dirname)   # ...Murphy and her ironclad Law
        root_dirname = root_dirname.rstrip(os.path.sep) + os.path.sep
        for pathname_part in pathname.split(os.path.sep):
            try:
                os.lstat(root_dirname + pathname_part)
            except OSError as exc:
                if hasattr(exc, 'winerror'):
                    if exc.winerror == ERROR_INVALID_NAME:
                        return False
                elif exc.errno in {errno.ENAMETOOLONG, errno.ERANGE}:
                    return False
    except TypeError as exc:
        return False
    else:
        return True


def map_deps_to_version(deps_arr: List[str]) -> Dict[str, str]:
    mapping = {}
    for d in deps_arr:
        if d.find(' => ') != -1:
            ds = d.split(' => ')
            d = ds[1]
        d = d.replace(" v", "[>v")  # might be able to just split on the empty space not _v and skip this :: insertion
        d_and_v = d.split("[>")
        mapping[d_and_v[0]] = d_and_v[1]
    return mapping


# argument checks
assert len(sys.argv) == 3, "need core repository and plugin repository path arguments"
core_repository_path = sys.argv[1]
plugin_repository_path = sys.argv[2]
assert is_pathname_valid(core_repository_path), "core repository path argument is not valid"
assert is_pathname_valid(plugin_repository_path), "plugin repository path argument is not valid"

# collect `go list -m all` output from both repositories; remain in the plugin repository
os.chdir(core_repository_path)
core_deps_b = subprocess.check_output(["go", "list", "-m", "all"])
os.chdir(plugin_repository_path)
plugin_deps_b = subprocess.check_output(["go", "list", "-m", "all"])
core_deps = core_deps_b.decode("utf-8")
core_deps_arr = core_deps.splitlines()
del core_deps_arr[0] # first line is the project repo itself
plugin_deps = plugin_deps_b.decode("utf-8")
plugin_deps_arr = plugin_deps.splitlines()
del plugin_deps_arr[0]
core_deps_mapping = map_deps_to_version(core_deps_arr)
plugin_deps_mapping = map_deps_to_version(plugin_deps_arr)

# iterate over dependency maps for both repos and find version conflicts
# attempt to resolve conflicts by adding adding a `require` for the core version to the plugin's go.mod file
none = True
for dep, core_version in core_deps_mapping.items():
    if dep in plugin_deps_mapping.keys():
        plugin_version = plugin_deps_mapping[dep]
        if core_version != plugin_version:
            print(f'{dep} has a conflict: core is using version {core_version} '
                  f'but the plugin is using version {plugin_version}')
            fixed_dep = f'{dep}@{core_version}'
            print(f'attempting fix by `go mod edit -require={fixed_dep}')
            subprocess.check_call(["go", "mod", "edit", f'-require={fixed_dep}'])
            none = False

if none:
    print("no conflicts to resolve")
    sys.exit(0)

# the above process does not work for all dep conflicts e.g. golang.org/x/text v0.3.0 will not stick this way
# so we will try the `go get {dep}` route for any remaining conflicts
updated_plugin_deps_b = subprocess.check_output(["go", "list", "-m", "all"])
updated_plugin_deps = updated_plugin_deps_b.decode("utf-8")
updated_plugin_deps_arr = updated_plugin_deps.splitlines()
del updated_plugin_deps_arr[0]
updated_plugin_deps_mapping = map_deps_to_version(updated_plugin_deps_arr)
none = True
for dep, core_version in core_deps_mapping.items():
    if dep in updated_plugin_deps_mapping.keys():
        updated_plugin_version = updated_plugin_deps_mapping[dep]
        if core_version != updated_plugin_version:
            print(f'{dep} still has a conflict: core is using version {core_version} '
                  f'but the plugin is using version {updated_plugin_version}')
            fixed_dep = f'{dep}@{core_version}'
            print(f'attempting fix by `go get {fixed_dep}')
            subprocess.check_call(["go", "get", fixed_dep])
            none = False

if none:
    print("all conflicts have been resolved")
    sys.exit(0)

# iterate over plugins `go list -m all` output one more time and inform whether or not the above has worked
final_plugin_deps_b = subprocess.check_output(["go", "list", "-m", "all"])
final_plugin_deps = final_plugin_deps_b.decode("utf-8")
final_plugin_deps_arr = final_plugin_deps.splitlines()
del final_plugin_deps_arr[0]
final_plugin_deps_mapping = map_deps_to_version(final_plugin_deps_arr)
none = True
for dep, core_version in core_deps_mapping.items():
    if dep in final_plugin_deps_mapping.keys():
        final_plugin_version = final_plugin_deps_mapping[dep]
        if core_version != final_plugin_version:
            print(f'{dep} STILL has a conflict: core is using version {core_version} '
                  f'but the plugin is using version {final_plugin_version}')
            none = False

if none:
    print("all conflicts have been resolved")
    sys.exit(0)

print("failed to resolve all conflicts")
sys.exit(1)