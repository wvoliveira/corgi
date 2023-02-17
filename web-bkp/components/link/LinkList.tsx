import React from "react";

export default function LinkList({ links }:{links:any}) {

  return (
    <>

      {links?.data?.map((link: any) => 
        <p key={link.id}>{link.id.substring(0, 10)}... {'=>'} {link.domain}/{link.keyword} {'=>'} {link.url}</p>
      )}

    </>
  )
}