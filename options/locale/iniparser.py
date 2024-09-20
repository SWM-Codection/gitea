#!/usr/bin/python3 
# -*- coding: utf-8 -*-

from configparser import ConfigParser, ExtendedInterpolation
from os import path 
from pathlib import Path

SUP_PATH = 'locale_en-US.ini'
SUB_PATH = 'locale_ko-KR.ini'
OUT_PATH = 'locale_ko-KR.mod.ini'

def to_relpath(x: str) -> str:
    return path.join(Path(__file__).resolve().parent, x)

def load_config(path: str) -> ConfigParser: 
    reader = ConfigParser(interpolation=None)
    with open(path, 'r') as f: 
        buf = '[__empty__]\n' + f.read()
        reader.read_string(buf)
    return reader

def save_config(path: str, obj: dict[str,object]) -> ConfigParser:
    writer = ConfigParser(interpolation=None) 
    writer.read_dict(obj)

    with open(path, 'w') as f:
        writer.write(f)

if __name__ == '__main__':
    sup = load_config(to_relpath(SUP_PATH))
    sub = load_config(to_relpath(SUB_PATH))
    common = {} 

    for section in sup.sections():
        common[section] = {}
        if section not in sub.sections(): continue
        for key in sup[section].keys():
            common[section][key] = sup[section][key]
            if key not in sub[section].keys(): continue
            common[section][key] = sub[section][key]

    save_config(to_relpath(OUT_PATH), common)
    
