import { TextInput, ActionIcon, Checkbox, Button, Group, Box } from '@mantine/core';
import { useForm } from '@mantine/form';
import { ChevronRight } from 'tabler-icons-react';

export function LinkInput() {

    const form = useForm({
      initialValues: {
        url: '',
      },
  
      validate: {
        url: (value) => (/^https:\/\/\S+$/.test(value) ? null : 'Invalid URL'),
      },
    });

  return (
    <Box sx={{ maxWidth: 300 }} mx="auto">
      <form onSubmit={form.onSubmit((values) => console.log(values))}>
        <Group position="left" mt="md" spacing='xs' >
          <TextInput
            required
            placeholder="https://www.google.com.br"
            {...form.getInputProps('url')}
          />

          <ActionIcon type="submit" variant="light" size="lg"><ChevronRight /></ActionIcon>
        </Group>
      </form>
    </Box>
  )
}