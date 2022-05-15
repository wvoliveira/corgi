import { TextInput, ActionIcon, Checkbox, Button, Group, Box, Alert } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useCallback, useState } from 'react';
import { AlertCircle, ChevronRight } from 'tabler-icons-react';
import APILink from '../lib/api/link';

export function LinkInput() {
  const [payload, setPayload] = useState({'url': ''});
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState({});
  const [error, setError] = useState(null);

  const sendRequest = useCallback(async (e) => {
    e.preventDefault();

    // Clear data and error.
    setData({});
    setError(null);

    // don't send again while we are sending
    if (loading) return
    // update state
    setLoading(true);
    
    console.log("Payload: " + payload.url);

    // Validate URL.
    if (!(/^(http|https):\/\/\S+/.test(payload.url))) {
      setError("Try a valid URL. Ex.: https://www.google.com.br")
      setLoading(false);
      return
    }

    // Send the request to create shortner link
    const response = await APILink.create(payload)
    if (response.data.status == "error") {
      setError(response.data.message);
      console.log(error);
    }

    if (response.data.status === "successful") {
      setData(response.data.data);
      console.log(data);
    }

    console.log("Response");
    console.log(response.data);

    // once the request is sent, update state again
    setLoading(false)
  }, [loading, payload]) // update the callback if the state changes

  return (
    <Box sx={{ maxWidth: 300 }} mx="auto">
      <form>
        <Group position="left" mt="md" spacing='xs' >
          <TextInput
            id="url"
            type="url"
            required
            error={error ? true : false}
            placeholder="https://"
            value={payload.url}

            onChange={(event) => {
              setPayload({"url": event.currentTarget.value})
              console.log(payload)
            }}
          />

          <ActionIcon 
            type="submit"
            variant="light" 
            size="lg" 
            loading={loading} 
            onClick={sendRequest}
          >
            <ChevronRight />
          </ActionIcon>

          {data.active ? <>
            <Alert icon={<AlertCircle size={16} />} title="Successful">
              http://{data.domain}/{data.keyword}
            </Alert>
          </> : null}

          {error ? <>
            <Alert icon={<AlertCircle size={16} />} title="Error" color="red">
              {error}
            </Alert>
          </> : null}
        </Group>
      </form>
    </Box>
  )
}