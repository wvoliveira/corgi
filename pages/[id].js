import { useRouter } from "next/router";
import redirect from 'nextjs-redirect'

import styles from "../styles/Home.module.css";
import { urls } from '../constants/urls.js'

function searchHash(hash) {
  const filtered = urls.filter((p) => p.hash === hash)

  var payload = {
    'id': '',
    'url': '',
    'status': 200,
    'message': 'successful',
  }

  if (filtered.length > 0) {
    const data = filtered[0]
    payload['id'] = data.id
    payload['url'] = data.url
  } else {
    payload['status'] = 404,
    payload['message'] = `url with id ${hash} not found.`
  }
  return payload
}

export default function URL() {
  const { query } = useRouter();
  const data = searchHash(query.id)

  const Redirect = redirect(data.url)

  return <div className={styles.container}>
    <div>
    {
      data ? 
        <Redirect />
      : <p className={styles.description}>loading...</p>
    }
    </div>
  </div>
}
