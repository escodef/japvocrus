from models.models import Translation, Sense
from bs4 import BeautifulSoup
import requests
from typing import List
from parsers.tatoeba_parser import TatoebaParser

class JardicParser():
    def __init__(self):
        super().__init__()
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:145.0) Gecko/20100101 Firefox/145.0'
        })
        self.tatoeba_parser = TatoebaParser()

    def parse_word(self, word: str) -> Translation:
        try:
            url = f"https://www.jardic.ru/search/search_r.php?q={word}&pg=0&sw=1472"
            
            response = self.session.get(url)
            response.raise_for_status()
            
            soup = BeautifulSoup(response.content, 'html.parser')
            
            word_element = soup.find('td', id='word-0-0')
            if not word_element:
                raise ValueError(f"Word element not found for: {word}")
            
            return self._parse_word_element(word_element, word)
            
        except Exception as e:
            self.logger.error(f"Failed to parse {word}: {e}")
            raise

    def _parse_word_element(self, word_element, original_word: str) -> Translation:
        word_ja = ""
        reading = ""
        
        word_spans = word_element.find_all('span')
        if word_spans:
            first_span = word_spans[0]
            word_ja = first_span.get_text(strip=True)
            
            if len(word_spans) > 1:
                reading_span = word_spans[1]
                reading = reading_span.get_text(strip=True)
        
        if not word_ja:
            font_element = word_element.find('font')
            if font_element:
                word_ja = font_element.get_text(strip=True)
        
        if not word_ja:
            word_ja = original_word
        
        senses = self._parse_senses(word_element)

        examples = self.tatoeba_parser.parse_word(original_word)
        
        return Translation(
            word=reading,
            reading=word_ja,
            senses=senses,
            examples=examples
        )

    def _parse_senses(self, word_element) -> List[Sense]:
        senses = []
        
        black_spans = word_element.find_all('span', style=lambda x: x and 'color: #000000' in x)
        
        for span in black_spans:
            text = span.get_text(separator='\n', strip=True)
            sense_lines = [line.strip() for line in text.split('\n') if line.strip()]
            
            for line in sense_lines:
                if line.startswith('•'):
                    sense_text = line[1:].strip()
                    
                    ru_text = sense_text
                    notes = ""
                    
                    if '(' in sense_text and ')' in sense_text:
                        start_idx = sense_text.find('(')
                        end_idx = sense_text.find(')', start_idx)
                        if start_idx != -1 and end_idx != -1:
                            ru_text = sense_text[:start_idx].strip()
                            notes = sense_text[start_idx + 1:end_idx].strip()
                    
                    sense = Sense(
                        ru=ru_text,
                        notes=notes,
                    )
                    senses.append(sense)
        
        if not senses:
            senses = self._parse_senses_alternative(word_element)
        
        return senses

    def _parse_senses_alternative(self, word_element) -> List[Sense]:
        """Нейронка говорит так тоже можно"""
        senses = []

        full_text = word_element.get_text(separator='\n', strip=True)
        lines = [line.strip() for line in full_text.split('\n') if line.strip()]
        
        for line in lines:
            if line.startswith('•'):
                sense_text = line[1:].strip()
                
                ru_text = sense_text
                notes = ""
                
                if '(' in sense_text and ')' in sense_text:
                    start_idx = sense_text.find('(')
                    end_idx = sense_text.find(')', start_idx)
                    if start_idx != -1 and end_idx != -1:
                        ru_text = sense_text[:start_idx].strip()
                        notes = sense_text[start_idx + 1:end_idx].strip()
                
                sense = Sense(
                    ru=ru_text,
                    notes=notes,
                )
                senses.append(sense)
        
        return senses