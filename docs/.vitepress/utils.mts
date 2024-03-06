// Adapted from https://github.com/vuejs/vitepress/issues/1737#issuecomment-1372010207
import * as fs from 'fs'
import { resolve } from 'path'

type ScanOptions = {
  capitalize?: boolean,
  replacements?: Record<string, string>,
  textFormat: (text: string) => string
}

// Path where our docs are living.
const BASE_PATH = resolve(__dirname, '../src/')

/**
 * Capitalize a given text.
 * @param {string} text - text to capitalize
 * @returns {string}
 */
function capitalize(text: string) {
  return text[0].toUpperCase() + text.slice(1)
}

/**
 * Given a folder path and an optional set of options, deep scan the folder for *.md files and return an array of items.
 * @param pathName The path to the folder with the *.md files to scan
 * @param options The options
 * @returns {Array}
 */
export function getSidebarItemsFromMdFiles(pathName: string, options: Partial<ScanOptions>) {
  const defaults: ScanOptions = {
    capitalize: true,
    textFormat: (text) => text.replace(/[-_]/g, ' ')
  }

	const path = resolve(BASE_PATH, `../src/${pathName}`)

	return getItems(path, {
    ...defaults,
    ...options
  })
}

// Read the folder and return the `{ text, items }` array.
function getItems(path: string, options: Partial<ScanOptions>) {
	let content = fs.readdirSync(path).filter(item => !item.startsWith('.'))

	if (!content) {
    return;
  }

  const getFormattedText = (text: string) => {
    let formattedText = options.textFormat(text)

    // If a custom label was provided, use that as-is and don't capitalize it to respect custom values.
    if (options.replacements && options.replacements[text]) {
      formattedText = options.replacements[text]
    } else
    // Otherwise, just check the `capitalize` flag.
    if (options.capitalize) {
      formattedText = capitalize(formattedText)
    }

    return formattedText;
  }

  let arr = content.map((item) => {
    const text = getFormattedText(item.split('.')[0])

    // If is content, just resolve it.
    if (item.endsWith('.md')) {
      return {
        text,
        link: resolve(path.replace(BASE_PATH, ''), item),
      }
    } else {
      const newPath = resolve(path, item)
      const isFolder = fs.lstatSync(newPath).isDirectory()

      if (!isFolder) {
        return
      }

      // If is a folder, recursively handle it.
      return {
        text,
        items: getItems(newPath, options),
        collapsible: true,
      }
    }
  })

  arr = arr.flatMap((item) => {
    if (item?.link) {
      item.link = normalizeLink(item.link)
    }

    return item
  })

  return arr
}

// Normalize a given path and return the `link` value in a standard, normalized format.
function normalizeLink(path) {
	return path.replace(/\\/g, '/').split('.')[0]
}
