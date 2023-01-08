const fetcher = (...args: any) => {
  console.log(...args);

  fetch(...args).
  then((res) => res.json());
}

export default fetcher;