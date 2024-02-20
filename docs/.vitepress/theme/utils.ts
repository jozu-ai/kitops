import { inBrowser } from 'vitepress'

// Get the current user OS
export function getUserOS(): 'mac' | 'windows' | 'linux' {
  // client-code only
  if (!inBrowser) {
    return
  }

  const osName = navigator.userAgent.toLowerCase();

  if (osName.indexOf('mac') !== -1) {
    return 'mac';
  }

  if (osName.indexOf('win') !== -1) {
    return 'windows';
  }

  return 'linux';
}
