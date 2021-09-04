import { useRouter } from "next/router";
import useSWR from "swr";
import redirect from 'nextjs-redirect'

import styles from "../styles/Home.module.css";

async function fetcher(...args) {
  const res = await fetch(...args)
  return res.json()
}

export default function URL() {
  const router = useRouter();
  const { query } = useRouter();

  console.log(query.id)

  const { data } = useSWR(`/api/hash/${query.id}`, fetcher)

  console.log(data)

  return <div style={{ textAlign: 'center' }}>
    <div>
    {
      data ? 
        window.location.href = data.url
      : 'loading...'
    }
    </div>
  </div>

}
