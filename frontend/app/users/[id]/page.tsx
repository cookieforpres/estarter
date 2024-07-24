import getRequestClient from "@/app/lib/getRequestClient"
import { APIError } from "@/app/lib/client"

export default async function UserDetails({
  params,
}: {
  params: { id: string }
}) {
  const client = getRequestClient()
  let response = null
  try {
    response = await client.user.GetUser(params.id)
  } catch (error) {
    console.log(error)
  }

  return (
    <section>
      {response !== null ? (
        <>
          <h1>Detail for {response.user.name}</h1>
          <p>ID: {response.user.id}</p>
        </>
      ) : (
        <p>You need to login to view this data</p>
      )}
    </section>
  )
}
