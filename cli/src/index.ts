import chalk from "chalk";
import { Command } from "commander";
import { getSdk, GetRecipeByUuidQuery, Recipe } from "./generated/graphql";
import { GraphQLClient } from "graphql-request";
import YAML from "yaml";

const program = new Command();

const client = new GraphQLClient("http://localhost:4242/query");
const sdk = getSdk(client);

async function getRecipe({ uuid }: { uuid: string }) {
  console.log(chalk.blue(`getting recipe ${uuid}`));

  const recipe = await sdk.getRecipeByUUID({ uuid });

  console.log(`GraphQL data:`, recipe);
  removeMeta(recipe);
  console.log(YAML.stringify(recipe));
}
async function main() {
  program
    .version("0.0.1")
    .description("CLI for food");

  program
    .command("get")
    .option("-U, --uuid <abc>", "specify the uuid")
    .action(getRecipe);
  await program.parseAsync(process.argv);

  if (!program.args.length) program.help();
}

main();

function removeMeta(obj: Record<string, string | number | any>) {
  for (let prop in obj) {
    if (
      prop === "uuid" ||
      prop === "__typename" ||
      obj[prop] === "" ||
      obj[prop] === 0 ||
      Object.keys(obj[prop]).length === 0
    )
      delete obj[prop];
    else if (typeof obj[prop] === "object") removeMeta(obj[prop]);
  }
}
