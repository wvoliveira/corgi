import Router from "next/router";
import {LoginInputs} from "../pages/auth/login";
import {catchAxiosError} from "./error";
import {post} from "./rest_service";

export async function login(inputs: LoginInputs): Promise<string | void> {
    const data = {"email": inputs.email, "password": inputs.password};
    const res: any = await post("/api/auth/password/login", data).catch(catchAxiosError);
    if (res.error) {
        return res.error;
    }
    await Router.push("/");
}