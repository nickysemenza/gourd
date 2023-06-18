import { EntitySelector } from "../../components/EntitySelector";
import { useNavigate } from "react-router-dom";
import { RecipesApi } from "../../api/openapi-fetch";
import { getOpenapiFetchConfig } from "../../util/config";
import { SubmitHandler, useForm } from "react-hook-form";
import { Button } from "../../components/ui/Button";

const CreateRecipe: React.FC = () => {
  const history = useNavigate();

  type Inputs = {
    url: string;
  };

  const { register, handleSubmit } = useForm<Inputs>();
  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    const bar = new RecipesApi(getOpenapiFetchConfig());
    const recipe = await bar.scrapeRecipe({
      scrapeRecipeRequest: { url: data.url },
    });

    console.log({ recipe });
    history(`/recipe/${recipe.detail.id}`);
  };

  return (
    <div className="w-1/2">
      <EntitySelector
        tall
        createKind="recipe"
        showKind={["recipe"]}
        onChange={(a) => {
          history(`/recipe/${a.value}`);
          console.log(a);
        }}
      />
      <form onSubmit={handleSubmit(onSubmit)}>
        <input
          type="url"
          className="border-2 border-gray-300 w-64"
          {...register("url", { required: true })}
        />
        <Button type="submit">scrape</Button>
      </form>
    </div>
  );
};

export default CreateRecipe;