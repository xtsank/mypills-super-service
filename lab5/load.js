import http from 'k6/http';
import exec from 'k6/execution';
import { sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

const base = __ENV.BASE_URL || 'http://nginx';
const medicine = '00000000-0000-4000-8000-000000000042';
const warmup = Number(__ENV.WARMUP_SECONDS || 30);
const duration = Number(__ENV.MEASURE_SECONDS || 180);
const writerPort = Number(__ENV.APP_PORT || 8081);
const requests = new Counter('web5_requests');
const errors = new Rate('web5_errors');
const readLatency = new Trend('web5_read_ms', true);
const readBackend = new Counter('web5_read_backend');

function scenario(execName, rate, seconds, start, phase) {
  return {
    executor: 'constant-arrival-rate', exec: execName, rate, timeUnit: '1s',
    duration: `${seconds}s`, startTime: `${start}s`, gracefulStop: '3s',
    preAllocatedVUs: execName === 'read' ? 12 : 6,
    maxVUs: execName === 'read' ? 24 : 12,
    tags: { phase },
  };
}

export const options = {
  setupTimeout: '90s',
  summaryTrendStats: ['avg', 'med', 'p(95)', 'max'],
  scenarios: {
    warmup_read: scenario('read', 16, warmup, 0, 'warmup'),
    warmup_write: scenario('write', 4, warmup, 0, 'warmup'),
    measure_read: scenario('read', 16, duration, warmup + 3, 'measure'),
    measure_write: scenario('write', 4, duration, warmup + 3, 'measure'),
  },
  thresholds: {
    dropped_iterations: ['count==0'],
    'web5_errors{phase:measure}': ['rate<0.01'],
    'web5_read_ms{phase:measure,ingress:root}': ['p(95)<=1000'],
    'web5_read_ms{phase:measure,ingress:mirror}': ['p(95)<=1000'],
  },
};

function params(token, name) {
  const headers = { 'Content-Type': 'application/json' };
  if (token) headers.Authorization = `Bearer ${token}`;
  return { headers, tags: { name }, timeout: '3s' };
}

function create(path, payload, token) {
  const r = http.post(base + path, JSON.stringify(payload), params(token, path));
  if (r.status !== 201) throw new Error(`Setup failed for ${path}: HTTP ${r.status}`);
  return r.json();
}

export function setup() {
  const login = `weblab5_${Date.now()}`;
  const { token } = create('/api/v2/users', {
    login, email: `${login}@example.test`, password: `${login}_pass`,
    age: 30, weight: 70, is_pregnant: false, is_driver: false, allergies: [], illnesses: [],
  });
  const readCabinet = create('/api/v2/cabinets', { name: 'Чтение' }, token).id;
  const writeCabinet = create('/api/v2/cabinets', { name: 'Запись' }, token).id;
  const writeItems = [];
  for (const cabinet of [readCabinet, writeCabinet]) {
    for (let i = 1; i <= 20; i++) {
      const item = create(`/api/v2/cabinets/${cabinet}/items`, {
        medicine_id: medicine, date_of_manufacture: `2026-01-${String(i).padStart(2, '0')}`,
        quantity: 3, minimum_quantity: 1,
      }, token);
      if (cabinet === writeCabinet) writeItems.push(item.id);
    }
  }

  const readPath = `/api/v2/cabinets/${readCabinet}/items?limit=20`;
  const replicaPorts = new Set();
  for (let i = 0; i < 40 && replicaPorts.size < 2; i++) {
    const response = http.get(base + '/mirror' + readPath, params(token, 'setup_replica'));
    if (hasSeedItems(response)) {
      const port = finalPort(response);
      if (port === 8085 || port === 8086) replicaPorts.add(port);
    }
    if (replicaPorts.size < 2) sleep(0.1);
  }
  if (replicaPorts.size !== 2) throw new Error('Setup: both replica readers must see the seed data.');
  return { token, readPath, writeCabinet, writeItems };
}

function hasSeedItems(response) {
  if (response.status !== 200) return false;
  try {
    return response.json('total') === 20;
  } catch (_) {
    return false;
  }
}

function finalPort(response) {
  const address = response.headers['X-Upstream-Addr'] || '';
  const match = address.match(/:(\d+)$/);
  return match ? Number(match[1]) : 0;
}

function tagsFor(iteration) {
  return {
    phase: exec.scenario.name.startsWith('measure_') ? 'measure' : 'warmup',
    ingress: iteration % 2 === 0 ? 'root' : 'mirror',
  };
}

export function read(data) {
  const tags = tagsFor(exec.scenario.iterationInTest);
  const prefix = tags.ingress === 'mirror' ? '/mirror' : '';
  const r = http.get(base + prefix + data.readPath, {
    ...params(data.token, 'GET cabinet items'), responseType: 'text',
  });
  const port = finalPort(r);
  const peers = tags.ingress === 'mirror' ? [8084, 8085, 8086] : [writerPort, 8082, 8083];
  const ok = hasSeedItems(r) && peers.includes(port);
  requests.add(1, { ...tags, operation: 'read' });
  errors.add(!ok, { ...tags, operation: 'read' });
  if (ok) {
    readLatency.add(r.timings.duration, tags);
    readBackend.add(1, { ...tags, port: String(port) });
  }
}

export function write(data) {
  const iteration = exec.scenario.iterationInTest;
  const tags = tagsFor(iteration);
  const prefix = tags.ingress === 'mirror' ? '/mirror' : '';
  const id = data.writeItems[iteration % data.writeItems.length];
  const path = `${prefix}/api/v2/cabinets/${data.writeCabinet}/items/${id}`;
  const quantity = 1 + (Math.floor(iteration / data.writeItems.length) % 5);
  const r = http.patch(base + path, JSON.stringify({ quantity }), {
    ...params(data.token, 'PATCH cabinet item'), responseType: 'none',
  });
  requests.add(1, { ...tags, operation: 'write' });
  errors.add(r.status !== 200 || finalPort(r) !== (prefix ? 8084 : writerPort),
    { ...tags, operation: 'write' });
}
