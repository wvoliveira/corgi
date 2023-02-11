// this is the logger for the browser
import pino from 'pino'

const config = {
  serverUrl: '/api/log',
}

const pinoConfig = {
  browser: {
    asObject: true
  }
}

if (config.serverUrl) {
  pinoConfig.browser.transmit = {
    level: 'info',
    send: (level, logEvent) => {
      const msg = logEvent.messages[0]

      let payload = JSON.stringify({
        time: logEvent.ts,
        message: msg,
        level,
        from: "app"
      })

      const headers = {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Headers': 'Origin, X-Requested-With, Content-Type, Accept',
        type: 'application/json'
      }

      let blob = new Blob([payload], headers)
      navigator.sendBeacon(`${config.serverUrl}`, blob)
    }
  }
}

const logger = pino(pinoConfig)

export const log = msg => logger.info(msg)
export default logger
