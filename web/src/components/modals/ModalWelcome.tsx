import Typewriter from "typewriter-effect";

export default function ModalWelcome() {
  return (
    <div className="fixed w-1/5 bg-white border rounded-lg left-6 font-mono p-4 h-3/4">
      <img src="/hand.png" className="w-28" />
      <div className="font-bold text-2xl">Welcome to your dashoard</div>
      <div className="mt-4">
        <Typewriter
          onInit={(typewriter) => {
            typewriter
              .typeString(
                "Hi, This is where you can create and manage your projects."
              )
              .pauseFor(2500)
              .typeString(" Click on the + button to create a new project.")
              .pauseFor(4000)
              .typeString(
                " We have eraly load your local folder you can now add a other service or a other folder."
              )
              .start();
          }}
          options={{
            delay: 30,
          }}
        />
      </div>
    </div>
  );
}
