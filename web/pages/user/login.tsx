import Head from "next/head";
import React from "react";

import CustomLink from "../../components/common/CustomLink";
import LoginForm from "../../components/profile/LoginForm";

const Login = () => (
  <>
    <Head>
      <title>Login | Corgi</title>
      <meta 
        name="description" 
        content="Please login to use fully-featured Corgi site. (Create custom links, groups, change your profile, etc.)" 
      />
    </Head>
    <div>
      <div>
        <div>
          <div>
            <h1>Login</h1>
            <p>
              <CustomLink href="/user/register" as="/user/register">
                Need an account?
              </CustomLink>
            </p>
            <LoginForm />
          </div>
        </div>
      </div>
    </div>
  </>
);

export default Login;
