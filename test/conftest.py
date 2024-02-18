import pytest
import requests
from urllib.parse import urljoin

class MonCronApi:
    def __init__(self, url: str):
        self._url = url
    
    def request(self, method: str, path: str, **kwargs) -> requests.Response:
        return getattr(requests, method)(urljoin(self._url, path), **kwargs)


def pytest_addoption(parser):
    parser.addoption("--url", action="store", required=True)

@pytest.fixture(scope="session")
def api(pytestconfig):
    return MonCronApi(pytestconfig.getoption("url"))

