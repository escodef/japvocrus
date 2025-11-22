import csv
import os
import logging
import time
from typing import List
from models.models import DictionaryType
from parsers.jardic_parser import JardicParser
from utils.file_utils import save_dictionary
from dotenv import load_dotenv
import random

load_dotenv()

csv_name = os.getenv("CSV_FILE")
filter = ["助動詞", "記号", "動詞-接尾"]

logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('logs/parser.log', encoding='utf-8'),
        logging.StreamHandler()
    ]
)

class JapaneseDictionaryParser:
    def __init__(self):
        self.parser = JardicParser()
        self.dictionary: DictionaryType = {}
    
    def parse_words(self, words: List[List[str]]) -> DictionaryType:
        """Парсит список слов и возвращает словарь"""
        for word in words:
            try:
                logging.info(f"Parsing word: {word[0]}")
                
                translation = self.parser.parse_word(word)
                if translation is None: 
                    continue

                self.dictionary[word[0]] = translation
                time.sleep(random.randint(1, 3))
                
            except Exception as e:
                logging.error(f"Failed to parse {word}: {e}")
                continue
        
        save_dictionary(self.dictionary, './dictionary.json')
        return self.dictionary

def main():
    words_to_parse = get_words()
    
    parser = JapaneseDictionaryParser()
    dictionary = parser.parse_words(words_to_parse)
    
    logging.info(f"Successfully parsed {len(dictionary)} words")

def get_words() -> List[List[str]]:
    a = []
    with open(csv_name, 'r', newline='', encoding='utf-8') as f:
        csv_file = csv.reader(f)
        for row in csv_file:
            if row[1] in filter or '【' in row[0]:
                continue
            a.append(row)
        return a

if __name__ == "__main__":
    main()