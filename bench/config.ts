import "./deps.ts";

const getEnv = (key: string): string => {
  const value = Deno.env.get(key);
  if (value === undefined) throw Error(`Undefined environment variable ${key}`);
  else return value;
};

const broker = getEnv("BROKER");
const topic = getEnv("TOPIC");
const username = getEnv("SASL_USERNAME");
const password = getEnv("SASL_PASSWORD");

export { broker, password, topic, username };
