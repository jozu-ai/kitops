import cheerio from 'cheerio'
import fs from 'node:fs'

export type Post = {
  title: string,
  description: string,
  url: string,
  published_time: string,
  site_name: string,
  image: string,
  icon: string,
  url: string
}

const postsUrls: string[] = require('./posts.json')

export default {
  async load() {
    const posts = postsUrls.map(async (url): Post[] => {
      try {
        const html = await (await fetch(url)).text()

        // Using cheerio to parse the html into actual dom nodes that we can interact.
        const $ = cheerio.load(html)

        // Tiny helper
        const getMetaTag = (name) => (
          $(`meta[name=${name}]`).attr("content") ||
          $(`meta[property="og:${name}"]`).attr("content") ||
          $(`meta[property="twitter${name}"]`).attr("content")
        )

        const title = getMetaTag('title') || $('title').text()
        const description = getMetaTag('description')
        const site_name = getMetaTag('site_name')
        const image = getMetaTag('image') || $('meta[property="og:image:url"]').attr('content')
        const icon = $('link[rel="icon"]').attr('href') || $('link[rel="shortcut icon"]').attr('href') || $('link[rel="alternate icon"]').attr('href')
        const author = getMetaTag('author')
        const published_time = $('meta[property="article:published_time"]').attr('content')

        return {
          url,
          title,
          description,
          published_time,
          site_name,
          image,
          icon
        } as Post
      }
      catch (e) {
        console.error(e)
      }
    })

    return Promise.all(posts)
  }
}
