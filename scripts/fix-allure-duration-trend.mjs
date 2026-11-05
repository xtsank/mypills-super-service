import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';

function correctDurationTrend(history, runTrend, durationTrend) {
  if (runTrend.length !== durationTrend.length) {
    throw new Error('Allure trend files have different run counts');
  }

  const tests = Object.values(history)
    .flatMap(({ items }) => items.map(({ time }) => time))
    .sort((left, right) => right.start - left.start);
  let offset = 0;
  let corrected = 0;

  const result = durationTrend.map((entry, index) => {
    const count = runTrend[index].data.total;
    const runTests = tests.slice(offset, offset + count);
    offset += count;
    if (runTests.length !== count || count === 0) {
      throw new Error(`Cannot recover Allure run ${index}: missing test history`);
    }

    const elapsed = Math.max(...runTests.map(({ stop }) => stop))
      - Math.min(...runTests.map(({ start }) => start));
    const testDuration = runTests.reduce((sum, { duration }) => sum + duration, 0);
    const oldDuration = entry.data.duration;
    if (oldDuration !== elapsed && oldDuration !== testDuration) {
      throw new Error(`Cannot match Allure run ${index}: recorded ${oldDuration} ms, expected ${elapsed} ms`);
    }
    if (oldDuration !== testDuration) corrected++;
    return { ...entry, data: { ...entry.data, duration: testDuration } };
  });

  return { result, corrected };
}

function jsonFile(path) {
  return JSON.parse(readFileSync(path, 'utf8'));
}

function correctDirectory(directory, historyDirectory = directory) {
  const history = jsonFile(join(historyDirectory, 'history.json'));
  const runTrend = jsonFile(join(directory, 'history-trend.json'));
  const durationPath = join(directory, 'duration-trend.json');
  const { result, corrected } = correctDurationTrend(history, runTrend, jsonFile(durationPath));
  if (corrected) writeFileSync(durationPath, JSON.stringify(result));
  console.log(`${durationPath}: ${corrected} of ${result.length} runs corrected`);
}

function correctReport(directory) {
  const historyDirectory = join(directory, 'history');
  correctDirectory(historyDirectory);
  correctDirectory(join(directory, 'widgets'), historyDirectory);
}

function correctSingleFile(path) {
  let html = readFileSync(path, 'utf8');
  const embedded = (name) => {
    const expression = new RegExp(`d\\('${name.replaceAll('.', '\\.')}','([A-Za-z0-9+/=]+)'\\)`);
    const match = html.match(expression);
    if (!match) throw new Error(`Allure single-file report has no ${name}`);
    return { match: match[0], data: JSON.parse(Buffer.from(match[1], 'base64').toString('utf8')) };
  };

  const history = embedded('history/history.json');
  let totalCorrected = 0;
  for (const directory of ['history', 'widgets']) {
    const runTrend = embedded(`${directory}/history-trend.json`);
    const durationTrend = embedded(`${directory}/duration-trend.json`);
    const { result, corrected } = correctDurationTrend(history.data, runTrend.data, durationTrend.data);
    if (corrected) {
      const replacement = `d('${directory}/duration-trend.json','${Buffer.from(JSON.stringify(result)).toString('base64')}')`;
      html = html.replace(durationTrend.match, replacement);
    }
    totalCorrected += corrected;
  }
  if (totalCorrected) writeFileSync(path, html);
  console.log(`Allure single-file duration trend: ${totalCorrected} values corrected`);
}

const [mode, path] = process.argv.slice(2);
if (mode === 'history' && path) correctDirectory(path);
else if (mode === 'report' && path) correctReport(path);
else if (mode === 'single-file' && path) correctSingleFile(path);
else {
  console.error('Usage: node scripts/fix-allure-duration-trend.mjs {history DIRECTORY|report DIRECTORY|single-file HTML_FILE}');
  process.exitCode = 2;
}
