from urllib.parse import urljoin
from conftest import MonCronApi
from random import choice, randrange
from math import ceil
import pytest

def test_num_hosts(api: MonCronApi):
    r = api.request('get', '/api/hosts')
    assert r.status_code == 200
    resp = r.json()
    assert resp['Total'] == 50
    assert len(resp['Data']) == 50

@pytest.mark.parametrize("size", list(range(1,51)))
def test_host_pagination(api: MonCronApi, size: int):
    resp = api.request('get', '/api/hosts').json()
    received_total = resp['Total']
    actual_total = 0

    all_hosts = []

    valid_pages = ceil(received_total / size)
    for page in range(valid_pages):
        page_resp = api.request('get', '/api/hosts', params={'p': page, 'size': size}).json()
        assert len(page_resp['Data']) > 0
        actual_total += len(page_resp['Data'])
        all_hosts.extend(page_resp['Data'])
    assert actual_total == received_total
    assert len(all_hosts) == received_total

    assert len(set(h['Name'] for h in all_hosts)) == received_total

@pytest.mark.parametrize("size", list(range(1, 15)))
def test_job_pagination(api: MonCronApi, size: int):
    resp = api.request('get', '/api/jobs').json()
    received_total = resp['Total']
    actual_total = 0

    all_jobs = []

    valid_pages = ceil(received_total / size)
    for page in range(valid_pages):
        page_resp = api.request('get', '/api/jobs', params={'p': page, 'size': size}).json()
        assert len(page_resp['Data']) > 0
        actual_total += len(page_resp['Data'])
        all_jobs.extend(page_resp['Data'])
    assert actual_total == received_total
    assert len(all_jobs) == received_total

    assert len(set(j for j in all_jobs)) == received_total

def test_runs_by_host(api: MonCronApi):
    actual_total = 0
    for num in range(1, 51):
        resp = api.request('get', '/api/runs', params={'host': f'srv-{num:03}.acme.corp'}).json()
        actual_total += resp['Total']
        # assert resp['Total'] == 18

        assert len(set(r['Host']['Name'] for r in resp['Data'])) == 1
    assert actual_total == 1000

def test_runs_by_host_and_name(api: MonCronApi):
    resp = api.request('get', '/api/runs', params={'host': 'srv-001.acme.corp', 'job': 'generate-thumbnails', 'size': 50}).json()
    assert resp['Total'] == 4

    assert len(set(r['Name'] for r in resp['Data'])) == 1
    assert len(set(r['Host']['Name'] for r in resp['Data'])) == 1