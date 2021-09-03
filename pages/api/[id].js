import { url } from '../../data'

export default function urlHandler({ query: { id } }, res) {
  const filtered = url.filter((p) => p.id === id)

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
    res.status(200).json(payload)
  } else {
    payload['status'] = 404,
    payload['message'] = `url with id ${id} not found.`
    res.status(404).json(payload)
  }
}