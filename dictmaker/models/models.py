from typing import List, Dict
from pydantic import BaseModel

class Example(BaseModel):
    ja: str
    ro: str
    tr: str

class Sense(BaseModel):
    ru: str
    notes: str = ""

class Translation(BaseModel):
    word: str
    reading: str
    senses: List[Sense] = []
    examples: List[Example] = []


DictionaryType = Dict[str, Translation]