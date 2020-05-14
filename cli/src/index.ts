import chalk from "chalk";
import { Command } from "commander";
import {
  getSdk,
  GetRecipeByUuidQuery,
  Recipe,
  RecipeInput,
} from "./generated/graphql";
import { GraphQLClient, ClientError } from "graphql-request";
import yaml from "js-yaml";
import fs from "fs-extra";
import path from "path";

const program = new Command();

const client = new GraphQLClient("http://localhost:4242/query");
const sdk = getSdk(client);

async function getRecipe({ uuid, output }: { uuid: string; output?: string }) {
  console.log(chalk.blue(`getting recipe ${uuid}`));
  try {
    const recipe = (await sdk.getRecipeByUUID({ uuid })).recipe;
    if (!recipe) {
      return;
    }
    // console.log(`GraphQL data:`, recipe);
    removeMeta(recipe);
    recipe.uuid = uuid;
    const r2: RecipeInput = recipe;
    r2.sections?.forEach((section, x) =>
      section.ingredients.forEach((i, y) => {
        i.name = recipe.sections[x].ingredients[y].info.name;
        delete (i as any).info;
        // delete r2.sections[x].ingredients[y].info;
      })
    );
    const r = yaml.safeDump(recipe);
    if (output) {
      fs.outputFileSync(`${output}/${recipe.name}.yaml`, r, "utf8");
    } else {
      console.log(r);
    }
  } catch (err) {
    console.error(
      chalk.red(
        "failed to fetch recipe: ",
        JSON.stringify((<ClientError>err).response.errors)
      ),
      err
    );
  }
}

const loadSingle = async (file: string) => {
  const contents: RecipeInput = yaml.safeLoad(fs.readFileSync(file, "utf8"));
  console.log(JSON.stringify(contents));

  try {
    await sdk.updateRecipe({ recipe: contents });
  } catch (e) {
    console.error(e);
  }
};
async function load({ input }: { input?: string }) {
  if (!input) {
    return;
  }
  var data = fs.statSync(input);
  if (data.isFile()) {
    loadSingle(input);
  } else if (data.isDirectory()) {
    const filenames = fs.readdirSync(input);

    console.log("\nCurrent directory filenames:");
    filenames.forEach((file) => {
      loadSingle(path.join(input, file));
    });
  }
}

async function getAllRecipes() {}

async function main() {
  program
    .version("0.0.1")
    .description("CLI for food");

  program
    .command("get")
    .option("-U, --uuid <abc>", "specify the uuid")
    .option(
      "-o, --output [dir]",
      "output directory (dir/name.yaml), defaults to STDout"
    )
    .action(getRecipe);
  program
    .command("load")
    .option(
      "-i, --input [dir]",
      "input directory (dir/name.yaml), defaults to STDin"
    )
    .action(load);
  program.command("all").action(getAllRecipes);
  await program.parseAsync(process.argv);

  if (!program.args.length) program.help();
}

main();

function removeMeta(obj: Record<string, string | number | any>) {
  for (let prop in obj) {
    if (
      prop === "uuid" ||
      prop === "__typename"
      // obj[prop] === "" ||
      // obj[prop] === 0 ||
      // Object.keys(obj[prop]).length === 0
    )
      delete obj[prop];
    else if (typeof obj[prop] === "object") removeMeta(obj[prop]);
  }
}
