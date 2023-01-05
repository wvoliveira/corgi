import React from "react";

export default function LinkList({ links }) {

  return (
    <>

      {links?.data?.map((link) => 
        <p>{link.id.substring(0, 10)}... => {link.domain}/{link.keyword} => {link.url}</p>
      )}

    </>
  )
}
