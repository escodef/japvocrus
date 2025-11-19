import json
from models.models import DictionaryType, Translation

def save_dictionary(dictionary: DictionaryType, filepath: str):
    """Сохраняет словарь в JSON файл"""
    dict_data = {
        word: translation.model_dump(mode='python')
        for word, translation in dictionary.items()
    }
    
    with open(filepath, 'w', encoding='utf-8') as f:
        json.dump(dict_data, f, ensure_ascii=False, indent=2)

def load_dictionary(filepath: str) -> DictionaryType:
    """Загружает словарь из JSON файла"""
    with open(filepath, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    return {
        word: Translation(**translation_data)
        for word, translation_data in data.items()
    }