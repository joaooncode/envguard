#!/usr/bin/env node
'use strict';

const os = require('os');
const path = require('path');
const fs = require('fs');
const https = require('https');
const { spawnSync } = require('child_process');
const tar = require('tar');

const pkg = require('../package.json');
const VERSION = pkg.version;
const REPO = 'joaooncode/envguard';

const PLATFORM_MAP = { linux: 'linux', darwin: 'darwin', win32: 'windows' };
const ARCH_MAP = { x64: 'amd64', arm64: 'arm64' };

function resolveTarget() {
  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];
  if (!platform || !arch) {
    console.error(`envguard: unsupported platform/arch combination: ${process.platform}/${process.arch}`);
    process.exit(1);
  }
  return { platform, arch };
}

function binaryName(platform) {
  return platform === 'windows' ? 'envguard.exe' : 'envguard';
}

function cacheDir(platform, arch) {
  return path.join(os.homedir(), '.envguard', 'bin', VERSION, `${platform}_${arch}`);
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https
      .get(url, { headers: { 'User-Agent': 'envguard-npm' } }, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          file.close();
          fs.unlink(dest, () => {
            download(res.headers.location, dest).then(resolve, reject);
          });
          return;
        }
        if (res.statusCode !== 200) {
          file.close();
          reject(new Error(`request to ${url} failed with status ${res.statusCode}`));
          return;
        }
        res.pipe(file);
        file.on('finish', () => file.close(resolve));
      })
      .on('error', reject);
  });
}

async function ensureBinary() {
  const { platform, arch } = resolveTarget();
  const dir = cacheDir(platform, arch);
  const bin = path.join(dir, binaryName(platform));

  if (fs.existsSync(bin)) {
    return bin;
  }

  fs.mkdirSync(dir, { recursive: true });

  const asset = `envguard_${platform}_${arch}.tar.gz`;
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${asset}`;
  const archivePath = path.join(dir, asset);

  console.error(`envguard: downloading ${asset} (v${VERSION})...`);
  await download(url, archivePath);
  await tar.x({ file: archivePath, cwd: dir });
  fs.unlinkSync(archivePath);

  if (!fs.existsSync(bin)) {
    throw new Error(`envguard binary not found after extracting ${asset}`);
  }
  if (platform !== 'windows') {
    fs.chmodSync(bin, 0o755);
  }
  return bin;
}

(async () => {
  try {
    const bin = await ensureBinary();
    const result = spawnSync(bin, process.argv.slice(2), { stdio: 'inherit' });
    if (result.error) {
      throw result.error;
    }
    process.exit(result.status === null ? 1 : result.status);
  } catch (err) {
    console.error(`envguard: ${err.message}`);
    process.exit(1);
  }
})();
