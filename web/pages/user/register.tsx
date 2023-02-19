import Head from "next/head";
import React from "react";

import CustomLink from "../../components/common/CustomLink";
import RegisterForm from "../../components/profile/RegisterForm";

const Register = () => (
  <>
  <Head>
    <title>Register | Corgi</title>
    <meta name="description" content="Please register before login" />
  </Head>
  <div>
    <div>
      <div>
        <div>
          <h1>Register</h1>
          <p>
            <CustomLink href="/user/login" as="/user/login">
              Have an account?
            </CustomLink>
          </p>

          <RegisterForm />
        </div>
      </div>
    </div>
  </div>
  </>
);

export default Register;
