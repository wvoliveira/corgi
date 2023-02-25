import { useRouter } from "next/router";
import React, {useEffect} from "react";
// @ts-ignore
import useSWR, { mutate, trigger } from "swr";

import { SERVER_BASE_URL } from "../../lib/utils/constant";
import ArticleList from "../../components/article/ArticleList";
import CustomImage from "../../components/common/CustomImage";
import ErrorMessage from "../../components/common/ErrorMessage";
import Maybe from "../../components/common/Maybe";
import EditProfileButton from "../../components/profile/EditProfileButton";
import FollowUserButton from "../../components/profile/FollowUserButton";
import ProfileTab from "../../components/profile/ProfileTab";
import UserAPI from "../../lib/api/user";
import checkLogin from "../../lib/utils/checkLogin";
import fetcher from "../../lib/utils/fetcher";
import storage from "../../lib/utils/storage";

const Profile = () => {
  const router = useRouter();

  console.log("QUERY: ", router.query);
  const keyURL = `${SERVER_BASE_URL}/users/username/${encodeURIComponent(String(router.query?.pid))}`

  const { data, error } = useSWR(keyURL, fetcher);
  const { data: currentUser } = useSWR("user", storage);

  console.log("fetchedProfile: ", data)
  console.log("profileError: ", error)

  if (error) return <ErrorMessage message="Can't load profile" />;

  const profile = data?.data;
  console.log(profile);

  // const { username, bio, image, following } = profile;
  const username = profile?.username;
  console.log(username);

  console.log("currentUser", currentUser);

  const isLoggedIn = checkLogin(currentUser);
  const isUser = currentUser && username === currentUser?.username;

  console.log("isLoggedIn: ", isLoggedIn);
  console.log("isUser: ", isUser);

  const handleFollow = async () => {
    mutate(
      `${SERVER_BASE_URL}/users/${router.query?.pid}`,
      { profile: { ...profile, following: true } },
      false
    );
    UserAPI.follow(router.query?.pid);
    trigger(`${SERVER_BASE_URL}/users/${router.query?.pid}`);
  };

  const handleUnfollow = async () => {
    mutate(
      `${SERVER_BASE_URL}/users/${router.query?.pid}`,
      { profile: { ...profile, following: true } },
      true
    );
    UserAPI.unfollow(router.query?.pid);
    trigger(`${SERVER_BASE_URL}/users/${router.query?.pid}`);
  };

  return (
    <div>
      <div>
        <CustomImage
          // src={image}
          src="image here"
          alt="User's profile image"
          className="user-img"
        />
        <h4>{profile?.username}</h4>
        {/* <p>{bio}</p> */}
        <p>"bio here"</p>
        <EditProfileButton isUser={isUser} />
        <Maybe test={isLoggedIn}>
          <FollowUserButton
            isUser={isUser}
            username={profile?.username}
            // following={following}
            following="following here"
            follow={handleFollow}
            unfollow={handleUnfollow}
          />
        </Maybe>
      </div>

          <div>
            <div>
              <ProfileTab profile={profile} />
            </div>
            <ArticleList />
          </div>
    </div>
  );
};

// Profile.getInitialProps = async ({ query: { pid } }) => {
//   const { data: initialProfile } = await UserAPI.get(pid);
//   return { initialProfile };
// };

export default Profile;
