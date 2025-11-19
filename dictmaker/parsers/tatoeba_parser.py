from typing import List, Optional
from bs4 import BeautifulSoup
import requests
from models.models import Example


class TatoebaParser():
    def __init__(self):
        super().__init__()
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:145.0) Gecko/20100101 Firefox/145.0'
        })

    def parse_word(self, word: str) -> List[Example]:
        try:
            url = f"https://www.jardic.ru/search/search_r.php?q={word}&pg=0&dic_tatoeba=1&sw=1472"
            
            response = self.session.get(url, timeout=10)
            response.raise_for_status()
            
            soup = BeautifulSoup(response.content, 'html.parser')
            examples = []
            
            example_containers = soup.find_all('span', style=lambda x: x and 'color: #00007F' in x)
            
            for container in example_containers[:2]:
                example = self._parse_example_container(container, word)
                if example:
                    examples.append(example)
            
            return examples
            
        except Exception as e:
            self.logger.error(f"Failed to parse {word}: {e}")
            raise

    def _parse_example_container(self, container, original_word: str) -> Optional[Example]:
        try:
            japanese_text = self._extract_japanese_sentence(container, original_word)
            if not japanese_text:
                return None
            
            translation_text = self._extract_translation(container)
            if not translation_text:
                return None
            
            return Example(
                ja=japanese_text,
                ro="",
                tr=translation_text
            )
            
        except Exception as e:
            self.logger.error(f"Failed to parse example container: {e}")
            return None

    def _extract_japanese_sentence(self, container, original_word: str) -> str:
        """Извлекает полное японское предложение"""
        red_font = container.find('font', color="#BF0000")
        if not red_font:
            return ""
        
        full_text = ""
        for element in container.contents:
            if element.name == 'span' and element.get('style', '').find('color: #000000') != -1:
                break
            if hasattr(element, 'get_text'):
                full_text += element.get_text()
            elif isinstance(element, str):
                full_text += element
        
        return full_text.strip()

    def _extract_translation(self, container) -> str:
        """Извлекает перевод из черного span'а"""
        black_span = container.find('span', style=lambda x: x and 'color: #000000' in x)
        if not black_span:
            return ""
        
        translation_text = black_span.get_text(separator='\n', strip=True)
        lines = [line.strip() for line in translation_text.split('\n') if line.strip()]
        
        return lines[0] if lines else ""