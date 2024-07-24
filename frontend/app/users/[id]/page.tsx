import getRequestClient from "@/app/lib/getRequestClient"
import { APIError, ErrCode } from "@/app/lib/client"

export default async function UserDetails({
  params,
}: {
  params: { id: number }
}) {
  const client = getRequestClient()
  let response = null
  try {
    response = await client.users.Get(params.id)
  } catch (error) {
    let message = ''

    if (error instanceof APIError) {
      message = error.message
      if (message === 'invalid auth param') {
        message = 'you need to login to view this data'
      }
    } else {
      message = String(error)
    }

    return (
      <>
        <p>{message.charAt(0).toUpperCase() + message.slice(1) + '!'}</p>
      </>
    )
  }

  return (
    <section>
      <h1>Detail for {response.user.name}</h1>
      <p>ID: {response.user.id}</p>
    </section>
  )
}
