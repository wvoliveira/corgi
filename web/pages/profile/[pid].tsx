import { useRouter } from "next/router";
import React, {useEffect} from "react";
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
  const [pid, setPID] = React.useState("");
  const [initialProfile, setInitialProfile] = React.useState({});

  console.log("QUERY: ", router.query);

  useEffect( () => {
    // fetch('/api/profile-data')
    //     .then((res) => res.json())
    //     .then((data) => {
    //       setData(data)
    //       setLoading(false)
    //     })

    setPID(router.query.pid)

    fetch(`${SERVER_BASE_URL}/users/username/${encodeURIComponent(String(router.query.pid))}`)
        .then((res) => res.json())
        .then((data) => {
          setInitialProfile(data)
        })

  }, [router.query])

  console.log("initialProfile: ", initialProfile);

  const {
    data: fetchedProfile,
    error: profileError,
  } = useSWR(
    `${SERVER_BASE_URL}/users/username/${encodeURIComponent(String(pid))}`,
    fetcher,
      { initialData: initialProfile }
  );

  console.log("fetchedProfile: ", fetchedProfile)
  console.log("profileError: ", profileError)

  // return <></>;

  if (profileError) return <ErrorMessage message="Can't load profile" />;

  // const profile = fetchedProfile?.data;

  const { data } = fetchedProfile || initialProfile;
  console.log(data);


  // const { username, bio, image, following } = profile;
  const username = data?.username;
  console.log(username);

  const { data: currentUser } = useSWR("user", storage);
  console.log("currentUser", currentUser);
  return <></>;

  const isLoggedIn = checkLogin(currentUser);
  const isUser = currentUser && username === currentUser?.username;

  console.log("isLoggedIn: ", isLoggedIn);
  console.log("isUser: ", isUser);
  return <></>;

  const handleFollow = async () => {
    mutate(
      `${SERVER_BASE_URL}/users/${pid}`,
      { profile: { ...profile, following: true } },
      false
    );
    UserAPI.follow(pid);
    trigger(`${SERVER_BASE_URL}/users/${pid}`);
  };

  const handleUnfollow = async () => {
    mutate(
      `${SERVER_BASE_URL}/users/${pid}`,
      { profile: { ...profile, following: true } },
      true
    );
    UserAPI.unfollow(pid);
    trigger(`${SERVER_BASE_URL}/users/${pid}`);
  };

  return (
    <div className="profile-page">
      <div className="user-info">
        <div className="container">
          <div className="row">
            <div className="col-xs-12 col-md-10 offset-md-1">
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
          </div>
        </div>
      </div>

      <div className="container">
        <div className="row">
          <div className="col-xs-12 col-md-10 offset-md-1">
            <div className="articles-toggle">
              <ProfileTab profile={profile} />
            </div>
            <ArticleList />
          </div>
        </div>
      </div>
    </div>
  );
};

// Profile.getInitialProps = async ({ query: { pid } }) => {
//   const { data: initialProfile } = await UserAPI.get(pid);
//   return { initialProfile };
// };

export default Profile;
