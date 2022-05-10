import Head from 'next/head'
import { Box, Input, Center, Container } from '@mantine/core';
import { LinkInput } from '../components/LinkInput';

export default function Home() {
  return (
    <div>
      <Head>
        <title>Corgi</title>
        <meta name="description" content="Shortener app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>

          <Box sx={(theme) => ({
              padding: theme.spacing.xl,
              borderRadius: theme.radius.md,

          })}>
            <LinkInput />            
            </Box>
      </main>

      {/* <footer>
        <a
          href="https://elga.io"
          target="_blank"
          rel="noopener noreferrer"
        >
          Powered by ELGA
        </a>
      </footer> */}
    </div>
  )
}
